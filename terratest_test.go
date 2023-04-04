package terratestPlus

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var DeployedResources *Deployment

func TestEndToEndOfSetupAndDeploy(t *testing.T) {

	DeployedResources = new(Deployment)
	DeployedResources.SetupTerraform(t, nil)
	DeployedResources.TerraformOptions.Logger = logger.Default
	DeployedResources.DeployInfrastructure()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

	// defer are stacked Last In First Out, so this is the last one in, it will
	//be run first, allowing the defer require above that recover be nil to check
	//Cleanup as part of the E2E
	defer DeployedResources.Cleanup()

	assert.False(t, DeployedResources.performCleanup, "Cleanup flag was not set False")

}

func TestDestroyOfSetupAndDeploy(t *testing.T) {
	DeployedResources.TeardownTerraform()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

}
