package k8s

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mednax-it/terratest-plus/deployment"
)

// * K8sDeployment is an extension of the deployment.D basic struct that includes additional values such as kube file locations and helm/helmfile functionality
type K8sDeployment struct {
	deployment.D

	KubeFiles map[string]string
}

func (k8s *K8sDeployment) GetKubeFiles(kubeFileListOutputIdentifier string) {
	k8s.KubeFiles = terraform.OutputMap(k8s.T, &k8s.TerraformOptions, kubeFileListOutputIdentifier)
}
