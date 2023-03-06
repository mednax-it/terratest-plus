package terratestPlus

import (
	"os"
	"os/user"
	"reflect"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/mednax-it/terratest-plus/deployment"
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
}

func (d *Deployment) getTFSource(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_source_dir"); present {
		d.TerraformSourceDir = val
	} else {
		d.TerraformSourceDir = options.TerraformDirectoryPath
	}
}

func (d *Deployment) getTFVars(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_var_file"); present {
		d.VarFilePath = val
	} else {
		d.VarFilePath = options.VarFileDirectoryPath + "local.tfvars"
	}

	if strings.Contains(d.VarFilePath, "local") {
		d.ExecutingInLocal = true
		d.RunInit = true
	}

	terraform.GetAllVariablesFromVarFile(d.T, d.TerraformSourceDir+d.VarFilePath, &d.VarFileValues)
}

func (d *Deployment) getTFBackend(options *SetupTerraformOptions) {
	if val, present := os.LookupEnv("TF_backend"); present {
		d.BackendFilePath = val
	} else {
		d.BackendFilePath = options.BackendDirectoryPath + "config.test_backend.tfbackend"
	}
}

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

// /*
// SetLocalWorkspace will be called if the TF_VAR_git_sha or the circle context of CIRCLE_SHA1 variable is not already set (as it will in the pipeline).

// It will attempt first to get the workspace name from the `git_sha` variable in the tfvars file.

// Failing to find that it will use the username in local.

// *NOTE!!* If running terratest in a docker, the username will always be root - so the gitsha (pushed to 7 chars) will be rootxxx.
// If multiple devs do not set the git_sha variable in the local.tfvars it will cause them to overwrite each others remote states!.
// */
// func setLocalWorkspace() string {

// 	var workspace string
// 	if varFileWorkspace, ok := varFileValues["git_sha"]; ok {
// 		workspace = varFileWorkspace.(string)

// 	} else {
// 		user, _ := user.Current()
// 		workspace = strings.ReplaceAll(user.Name, " ", "")
// 	}

// 	return workspace
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

// /*
// GetState calls the `terraform state pull` command and retrieves the current state file.

// This should be run after Apply.

// This allows for a map for use in testing certain situations.

// GetState sets the State global var for use in various tests.

// GetState should not need to be used in General Testing but is provided exposed just in case. Instead use the helper functions found in wrapper/state.go.

// Note: in order to make matching to a Struct easier, the terraform module names (which are usually like `["name"]` ) have been cleaned to just `[name]`.
// */
// func GetState(t *testing.T) {
// 	tf_get_state, err := terraform.RunTerraformCommandAndGetStdoutE(t, &GeneralTerraformOptions, "state", "pull")

// 	tf_get_state = strings.ReplaceAll(tf_get_state, "\n", "")
// 	tf_get_state = strings.ReplaceAll(tf_get_state, "\\", "")
// 	tf_get_state = strings.ReplaceAll(tf_get_state, "[\"", "[")
// 	tf_get_state = strings.ReplaceAll(tf_get_state, "\"]", "]")

// 	if err != nil {
// 		assert.FailNow(t, "State File was not able to be pulled - tests cannot run.")
// 	}
// 	State = &TerraformState{}
// 	json.Unmarshal([]byte(tf_get_state), &RawState)
// 	_, err = marshmallow.Unmarshal([]byte(tf_get_state), State)

// 	if err != nil {
// 		assert.FailNow(t, "Could not build map of State File - Tests cannot run.")
// 	}

// }

// /*
// TeardownTerraform accepts a testing struct and the path to the directory where the terraform main.tf is to then Destroy.

// This respects path conventions, such as "../terraform/aks" to go back up one.
// */
// func TeardownTerraform(t *testing.T, terraformDirectory string) {

// 	// Clean up resources with "terraform destroy" at the end of the test.
// 	terraform.Destroy(t, &GeneralTerraformOptions)

// }
