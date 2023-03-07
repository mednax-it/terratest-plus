package terratestPlus

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"reflect"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/mednax-it/terratest-plus/deployment"
	"github.com/perimeterx/marshmallow"
	"github.com/stretchr/testify/assert"
)

type SetupTerraformOptions struct {
	TerraformDirectoryPath string `default:"src/"`
	VarFileDirectoryPath   string `default:"vars/"`
	BackendDirectoryPath   string `default:"backends/"`
	Workspace              string
}

func defaultValues(o *SetupTerraformOptions) {
	if o.TerraformDirectoryPath == "" {
		o.TerraformDirectoryPath = getDefaultTag(*o, "TerraformDirectoryPath")
	}

	if o.VarFileDirectoryPath == "" {
		o.VarFileDirectoryPath = getDefaultTag(*o, "VarFileDirectoryPath")
	}

	if o.BackendDirectoryPath == "" {
		o.BackendDirectoryPath = getDefaultTag(*o, "BackendDirectoryPath")
	}

	if o.Workspace == "" {
		user, _ := user.Current()
		o.Workspace = strings.ReplaceAll(user.Name, " ", "")
	}
}

/*
	GetDefaultTag returns the value in `default:"string"` tag on a struct.

Only works for strings, must be called directly per struct attribute
*/
func getDefaultTag(o SetupTerraformOptions, property string) string {
	typ := reflect.TypeOf(o)
	f, _ := typ.FieldByName(property)
	return f.Tag.Get("default")
}

type Deployment struct {
	deployment.D
}

/*
deployment.SetupTerraform setsup all the values needed for finding the terraform in git and various var/backends/workspaces.

You can pass nil to [options] and it will set the defaults noted in the Readme.md

You can set env variables that will be used for certain values. See Readme.md again.
*/
func (d *Deployment) SetupTerraform(t *testing.T, options *SetupTerraformOptions) {
	//save the Testing struct for use in various functions within
	d.T = t
	if options == nil {
		options = new(SetupTerraformOptions)
	}
	defaultValues(options)

	d.getTFSource(options)
	d.getTFVars(options)
	d.getTFBackend(options)
	//d.getTFWorkspace(options)

	d.TerraformOptions = *terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: d.TerraformSourceDir,
		VarFiles:     []string{d.VarFilePath},
		EnvVars:      map[string]string{"TF_PARAM_BACKEND_CONFIG_FILE": d.BackendFilePath},
		Parallelism:  10,
		Logger:       logger.Discard,
	})
}

/*
	DeployInfrastructure will perform Terraform Init and Terraform Apply with the options that were set in SetupTerraform

Init is controlled by both the RunInit flag and the Test Structure Stage variable SKIP_terraformin_init

Apply is controlled by the Test Structure Stage variable SKIP_terraform_apply
*/
func (d *Deployment) DeployInfrastructure() {

	test_structure.RunTestStage(d.T, "terraform_init", func() {
		if d.RunInit {
			terraform.Init(d.T, &d.TerraformOptions)
		}
	})

	test_structure.RunTestStage(d.T, "terraform_apply", func() {
		terraform.WorkspaceSelectOrNew(d.T, &d.TerraformOptions, d.WorkspaceName)
		terraform.Apply(d.T, &d.TerraformOptions)
	})

	d.GetState()
}

/*
GetState calls the `terraform state pull` command and retrieves the current state file.

This should be run after Apply.

This allows for a map for use in testing certain situations.

GetState sets the State struct var for use in various tests, as well as the RawState var.

GetState should not need to be used in General Testing but is provided exposed just in case. Instead use the helper functions found in wrapper/state.go.

Note: in order to make matching to a Struct easier, the terraform module names (which are usually like `["name"]` ) have been cleaned to just `[name]`.
*/
func (d *Deployment) GetState() {
	tf_get_state, err := terraform.RunTerraformCommandAndGetStdoutE(d.T, &d.TerraformOptions, "state", "pull")

	tf_get_state = strings.ReplaceAll(tf_get_state, "\n", "")
	tf_get_state = strings.ReplaceAll(tf_get_state, "\\", "")
	tf_get_state = strings.ReplaceAll(tf_get_state, "[\"", "[")
	tf_get_state = strings.ReplaceAll(tf_get_state, "\"]", "]")

	if err != nil {
		logger.Logf(d.T, "Error pulling State file: %s", errors.Unwrap(err))
		assert.FailNow(d.T, "State File was not able to be pulled - tests cannot run.")
	}
	json.Unmarshal([]byte(tf_get_state), &d.RawState)
	_, err = marshmallow.Unmarshal([]byte(tf_get_state), &d.State)

	if err != nil {
		logger.Logf(d.T, "Error building State Map: %s", errors.Unwrap(err))
		assert.FailNow(d.T, "Could not build map of State File - Tests cannot run.")
	}

}

/*
TeardownTerraform accepts a testing struct and the path to the directory where the terraform main.tf is to then Destroy.

This respects path conventions, such as "../terraform/aks" to go back up one.
*/
func (d *Deployment) TeardownTerraform() {

	// Clean up resources with "terraform destroy" at the end of the test.
	terraform.Destroy(d.T, &d.TerraformOptions)

}

/* GetTFSource checks the env variable TF_source_dir first, then uses the values from the passed in options
 */
func (d *Deployment) getTFSource(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_source_dir"); present {
		d.TerraformSourceDir = val
	} else {
		d.TerraformSourceDir = options.TerraformDirectoryPath
	}
}

/*
	GetTFVars looks for the env variable TF_var_file first then takes from the options.

# If the options contains a full path (ending in .tfvars) then it takes it as is, else it defaults to local.tfvars

If the Varfile name contains the word 'local' it also sets the ExecutingInLocal and RunInit bools to true.
*/
func (d *Deployment) getTFVars(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_var_file"); present {
		d.VarFilePath = val
	} else {
		if strings.Contains(options.VarFileDirectoryPath, ".tfvars") {
			d.VarFilePath = options.VarFileDirectoryPath
		} else {
			d.VarFilePath = options.VarFileDirectoryPath + "local.tfvars"
		}
	}

	if strings.Contains(d.VarFilePath, "local") {
		d.ExecutingInLocal = true
		d.RunInit = true
	}

	terraform.GetAllVariablesFromVarFile(d.T, d.TerraformSourceDir+d.VarFilePath, &d.VarFileValues)
}

/*
	GetTFBackend looks for the env variablle TF_backend first, then takes from the options.

If the passed in options contains the word ".tfbackend" as a full path then it is used as is.
Otherwise it defaults to `config.test_backend.tfbackend`
*/
func (d *Deployment) getTFBackend(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_backend"); present {
		d.BackendFilePath = val
	} else {
		if strings.Contains(options.BackendDirectoryPath, ".tfbackend") {
			d.BackendFilePath = options.BackendDirectoryPath
		} else {
			d.BackendFilePath = options.BackendDirectoryPath + "config.test_backend.tfbackend"
		}

	}
}

/*
	GetTFWorkspace looks for the env variable TF_workspace first, then the CIRCLE_SHA, then takes the passed workspace name.

It sets all of them to 7 characters, either cutting it down or adding 0s.
*/
func (d *Deployment) getTFWorkspace(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_workspace"); present {
		d.WorkspaceName = val
	} else if val, present := os.LookupEnv("CIRCLE_SHA1"); present {
		d.WorkspaceName = val

	} else {
		d.WorkspaceName = options.Workspace
	}

	d.cleanWorkspaceName()
}

/*
	CleanWorkspaceName shortens the name to 7 characters and sets the env var and some display variables.

We do this here in the terratest wrapper so we it is only done in one place, and that place is along side other similar operations.
*/
func (d *Deployment) cleanWorkspaceName() {
	if len(d.WorkspaceName) >= 7 {
		d.WorkspaceName = d.WorkspaceName[0:7]
	} else {
		d.WorkspaceName += strings.Repeat("0", 7-len(d.WorkspaceName))
	}

	//os.Setenv("TF_VAR_git_sha", d.WorkspaceName)
}

// /*
// DiscoverNumberOfPlatforms consumes a varfile and reads in the 'platforms' list.
// It sets the NumberOfPlatforms variable to that value.
// */
// func discoverNumberOfPlatforms(t *testing.T) {

// 	NumberOfPlatforms = len(varFileValues["platforms"].([]interface{}))

// }

// /*
// 	SetupTerraform accepts a testing struct and the path to the directory where the terraform main.tf is.

// It will run Init if in the local environment, but in the pipeline (assuming all env variables were correctly set) it will skip for time as that was run in an early CircleCIJob.

// It will switch to the workspace found in TF_VAR_git_sha or CIRCLE_SHA1, prioritizing TF_VAR_git_sha and trim/fill it to 7 characters.

// if it cannot find either of these values, it will look in the tfvars file for a git_sha variable.

// if it cannot find one there, it will default to the username of the workspace running this test.

// Note: this means that if you do not set git_sha in the tfvars file and you are using docker to run tests, they will be run in the `rootxxx` workspace, which will overwrite other peoples work who are doing the same.

// To prevent this, all thats needed is a unique git_sha variable in the tfvars file.

// This respects path conventions, such as "../terraform/aks" to go back up one directory.
// */
// func SetupTerraform(t *testing.T, terraformDirectory string) {

// 	if runInit { // this is only set True if the local.tfvars is used, so its entirely for local testing.
// 		terraform.Init(t, &GeneralTerraformOptions)
// 	}

// 	workspace := os.Getenv("TF_VAR_git_sha")
// 	sha := os.Getenv("CIRCLE_SHA1")
// 	if workspace == "" && sha == "" {

// 		workspace = setLocalWorkspace()

// 	}
// 	if sha != "" {
// 		workspace = sha
// 	}
// 	workspace = cleanWorkspaceName(workspace)

// 	// Run Terraform Apply - Can be skipped by setting the env var `SKIP_terraform_apply=true` in order to speed up local testing
// 	test_structure.RunTestStage(t, "terraform_apply", func() {
// 		terraform.WorkspaceSelectOrNew(t, &GeneralTerraformOptions, workspace)

// 		terraform.Apply(t, &GeneralTerraformOptions)
// 	})

// 	GetState(t)
// 	GetPlatformDetails(terraform.OutputMapOfObjects(t, &GeneralTerraformOptions, "platform_values"))

// 	KubeFiles = terraform.OutputMap(t, &GeneralTerraformOptions, "kube_files")
// 	SubscriptionId = terraform.Output(t, &GeneralTerraformOptions, "subscription_id")

// }

// /* GetPlatformDetails casts the map of maps of strings into a map of SinglePlatformDetails structs.
//  */
// func GetPlatformDetails(output map[string]interface{}) {
// 	for region, singlePlatform := range output {
// 		val := singlePlatform.(map[string]interface{})
// 		PlatformDetails[region] = SinglePlatformDetail{
// 			Context:                val["kube_context"].(string),
// 			ApplicationGatewayName: val["application_gateway_name"].(string),
// 			PublicIp:               val["public_ip"].(string),
// 			ResourceGroupName:      val["resource_group_name"].(string),
// 			IdentityId:             val["identity_id"].(string),
// 			IdentityClientId:       val["identity_client_id"].(string)}
// 	}
// }
