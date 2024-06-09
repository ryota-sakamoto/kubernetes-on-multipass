package provisioner

import (
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/multipass"
)

type Config struct {
	Name       string
	CPUs       string
	Memory     string
	Disk       string
	K8sVersion string
	Image      string
}

func CreateMaster(config Config) error {
	instance, err := multipass.GetInstance(config.Name)
	if err != nil {
		return err
	}
	if instance != nil {
		return nil
	}

	template, err := cloudinit.Generate(GetMasterTemplate(), map[string]string{
		"KubernetesVersion": config.K8sVersion,
		"Arch":              "amd64",
	})
	if err != nil {
		return err
	}

	return multipass.LaunchInstance(multipass.InstanceConfig{
		Name:   config.Name,
		CPUs:   config.CPUs,
		Memory: config.Memory,
		Disk:   config.Disk,
		Image:  config.Image,
	}, template)
}

func CreateWorker(config Config) error {
	instance, err := multipass.GetInstance(config.Name)
	if err != nil {
		return err
	}
	if instance != nil {
		return nil
	}

	template, err := cloudinit.Generate(GetWorkerTemplate(), map[string]string{
		"KubernetesVersion": config.K8sVersion,
		"Arch":              "amd64",
	})
	if err != nil {
		return err
	}

	return multipass.LaunchInstance(multipass.InstanceConfig{
		Name:   config.Name,
		CPUs:   config.CPUs,
		Memory: config.Memory,
		Disk:   config.Disk,
		Image:  config.Image,
	}, template)
}
