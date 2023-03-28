package deployment

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type D struct {
	T                *testing.T
	TerraformOptions terraform.Options
	SubscriptionId   string
	State            *TerraformState
	RawState         map[string]interface{}
	BackendValues    map[string]interface{}
	VarFileValues    map[string]interface{}
	OutputValues     map[string]interface{}

	// Reference variables for display and output primarily
	RunInit            bool
	ExecutingInLocal   bool
	WorkspaceName      string
	VarFilePath        string
	BackendFilePath    string
	TerraformSourceDir string
}
