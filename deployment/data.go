package deployment

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type D struct {
	t                testing.T
	TerraformOptions terraform.Options
	State            *TerraformState
	RawState         map[string]interface{}
	SubscriptionId   string

	VarFileValues map[string]interface{}

	RunInit          bool
	ExecutingInLocal bool
	WorkspaceName    string
	VarFilePath      string
}
