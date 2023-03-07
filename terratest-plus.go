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
	performCleanup bool
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
	d.getOutputValues()
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

/*
	Clean up is a function to be defered and catch the end state of the testing, to determine what needs to be cleaned up.

If a panic occurs, the infrastructure will be destroyed.

if the testing fails and it is NOT being executed in the local (ie local iterative testing) then the infrastructure will be destroyed.

If SKIP_terraform_destroy env variable is set, even if the the above qualifiers occur for a destruction, it will not do so.
*/
func (d *Deployment) Cleanup() {
	if r := recover(); r != nil {
		logger.Log(d.T, "\n\n\033[91m>>> Catastrophic Error (Panic!). <<<\033[0m\n\n")
		logger.Log(d.T, r)
		d.performCleanup = true
	}

	if d.T.Failed() && !d.ExecutingInLocal {
		logger.Log(d.T, "\n\n\033[91m>>> One or more Tests failed. <<<\033[0m\n\n")
		d.performCleanup = true
	}

	if d.ExecutingInLocal {
		logger.Logf(d.T, "\n\n\033[93m>>> Local Testing - Env Left in place. Use the following when finished: \033[0m\n\n\t$ terraform workspace select %s\n\t$ terraform destroy -var-file=%s\n\t$ terraform workspace select default\n\t$ terraform workspace delete %s\n\n", wrapper.WorkspaceName, wrapper.VarfilePath, wrapper.WorkspaceName)

	}

	test_structure.RunTestStage(d.T, "terraform_destroy", func() {
		if d.performCleanup {
			logger.Log(d.T, "\n\n>>> Cleaning up after failure in testing ... <<< \n\n")
			d.TeardownTerraform()
		}
	})
}

/* RunTests takes a map of test functions and runs them through go GoRoutines
 */
func (d *Deployment) RunTests(dispatch map[string]func(t *testing.T)) {

	for name, testCommand := range dispatch {
		d.T.Run(name, testCommand)
	}
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
	GetOutputValues stores the outputs of the Terraform in a map.

Must be called after terraform apply has been run.

As such it is called automatically as part of the DeployInfrastructure
*/
func (d *Deployment) getOutputValues() {
	d.OutputValues = terraform.OutputAll(d.T, &d.TerraformOptions)
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

}
