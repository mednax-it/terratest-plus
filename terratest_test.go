package terratestPlus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var DeployedResources *Deployment

func TestEndToEndOfSetupAndDeploy(t *testing.T) {

	DeployedResources = new(Deployment)
	DeployedResources.SetupTerraform(t, nil)
	DeployedResources.DeployInfrastructure()

	defer func() {
		require.Nil(t, recover(), "Panic Occurred, test will fail")
	}()

	// defer are stacked Last In First Out, so this is the last one in, it will
	//be run first, allowing the defer require above that recover be nil to check
	//Cleanup as part of the E2E
	defer DeployedResources.Cleanup()

	require.False(t, DeployedResources.T.Failed(), "There was a failure in an assertion")

}

func TestDestroyOfSetupAndDeploy(t *testing.T) {
	DeployedResources.TeardownTerraform()
}
