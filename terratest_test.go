package terratestPlus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* Mock Terraform deployment testing

Note! This is not mocking the terraform deployment but rather actually deploying a mock terraform file.

As such, the order of tests in this file is very important! The first test is an actual deployment of resources, and the last test needs to be a cleanup.

You can always deploy and clean up in the same test, but in order to prevent dozens of deployments making it take longer, then we put the tests in specific order to reuse existing deployments.

This is not ideal, but it works for the time being.
*/

var DeployedResources *Deployment

func TestEndToEndOfSetupAndDeploy(t *testing.T) {

	DeployedResources = new(Deployment)
	DeployedResources.SetupTerraform(t, nil)
	//Uncomment below if you need to see terraform logs to diagnose an issue. Otherwise they are suppressed
	//Alternatively use LOG_TERRAFORM=true env variable
	//DeployedResources.TerraformOptions.Logger = logger.Default
	DeployedResources.DeployInfrastructure()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

	// defer are stacked Last In First Out, so this is the last one in, it will
	//be run first, allowing the defer require above that recover be nil to check
	//Cleanup as part of the E2E

	defer DeployedResources.Cleanup()

	//using local.tfvars so this should trigger ExecutingInLocal to be true
	assert.True(t, DeployedResources.ExecutingInLocal, "ExecutingInLocal was not set true")

	assert.False(t, DeployedResources.performCleanup, "Cleanup flag was not set False")

}

func TestGetStateThrowsNoErrors(t *testing.T) {

	DeployedResources.GetState()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

}

func TestDestroyOfSetupAndDeploy(t *testing.T) {
	DeployedResources.TeardownTerraform()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

}
