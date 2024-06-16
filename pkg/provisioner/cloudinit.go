package provisioner

import (
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
)

func GetMasterTemplate(k8sVersion string, arch string, isRegisterNode bool) cloudinit.Config {
	result := baseTemplate

	result.WriteFiles = append(result.WriteFiles, masterTemplate.WriteFiles...)
	result.RunCmds = append(result.RunCmds, masterTemplate.RunCmds...)
	result.Vars = map[string]any{
		"KubernetesVersion": k8sVersion,
		"Arch":              arch,
		"RegisterNode":      isRegisterNode,
	}

	return result
}

func GetWorkerTemplate(k8sVersion string, arch string) cloudinit.Config {
	result := baseTemplate
	result.Vars = map[string]any{
		"KubernetesVersion": k8sVersion,
		"Arch":              arch,
	}

	return result
}
