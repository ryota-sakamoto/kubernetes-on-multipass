package provisioner

import (
	"fmt"
	"log/slog"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/multipass"
)

func LaunchInstance(clusterName string, config InstanceConfig, cloudinitConfig cloudinit.Config) (string, error) {
	name := config.Name
	if name == "" {
		name = GetRandomName()
	}

	instanceName := fmt.Sprintf("%s-%s", clusterName, name)

	instance, err := multipass.GetInstance(instanceName)
	if err != nil {
		return "", fmt.Errorf("failed to get instance: %w", err)
	}

	slog.Debug("get instance", slog.String("instanceName", instanceName), slog.Any("instance", instance))
	if instance != nil {
		return "", fmt.Errorf("instance already exists: %s", instanceName)
	}

	template, err := cloudinit.Generate(GetWorkerTemplate(), map[string]string{
		"KubernetesVersion": config.K8sVersion,
		"Arch":              "amd64",
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate cloud-init template: %w", err)
	}

	return instanceName, multipass.LaunchInstance(multipass.InstanceConfig{
		Name:   instanceName,
		CPUs:   config.CPUs,
		Memory: config.Memory,
		Disk:   config.Disk,
		Image:  config.Image,
	}, template)
}
