package provisioner

import (
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
)

func GetMasterTemplate() cloudinit.Config {
	result := baseTemplate

	result.WriteFiles = append(result.WriteFiles, masterTemplate.WriteFiles...)
	result.RunCmds = append(result.RunCmds, masterTemplate.RunCmds...)

	return result
}

func GetWorkerTemplate() cloudinit.Config {
	return baseTemplate
}
