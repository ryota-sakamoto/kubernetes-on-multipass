package provisioner

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/kubernetes"
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

func CreateMaster(clusterName string, config Config) error {
	slog.Debug("create master", slog.String("clusterName", clusterName), slog.Any("config", config))

	instanceName := fmt.Sprintf("%s-%s", clusterName, "master")
	instance, err := multipass.GetInstance(instanceName)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	slog.Debug("get instance", slog.String("instanceName", instanceName), slog.Any("instance", instance))
	if instance != nil {
		return nil
	}

	template, err := cloudinit.Generate(GetMasterTemplate(), map[string]string{
		"KubernetesVersion": config.K8sVersion,
		"Arch":              "amd64",
	})
	if err != nil {
		return fmt.Errorf("failed to generate cloud-init template: %w", err)
	}

	return multipass.LaunchInstance(multipass.InstanceConfig{
		Name:   instanceName,
		CPUs:   config.CPUs,
		Memory: config.Memory,
		Disk:   config.Disk,
		Image:  config.Image,
	}, template)
}

func CreateWorker(clusterName string, config Config) error {
	slog.Debug("create worker", slog.String("clusterName", clusterName), slog.Any("config", config))

	name := config.Name
	if name == "" {
		name = GetRandomName()
	}

	instanceName := fmt.Sprintf("%s-%s", clusterName, name)

	instance, err := multipass.GetInstance(instanceName)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	slog.Debug("get instance", slog.String("instanceName", instanceName), slog.Any("instance", instance))
	if instance != nil {
		return nil
	}

	template, err := cloudinit.Generate(GetWorkerTemplate(), map[string]string{
		"KubernetesVersion": config.K8sVersion,
		"Arch":              "amd64",
	})
	if err != nil {
		return fmt.Errorf("failed to generate cloud-init template: %w", err)
	}

	return multipass.LaunchInstance(multipass.InstanceConfig{
		Name:   instanceName,
		CPUs:   config.CPUs,
		Memory: config.Memory,
		Disk:   config.Disk,
		Image:  config.Image,
	}, template)
}

func GenerateKubeconfig(name string) error {
	slog.Debug("generate kubeconfig", slog.String("name", name))

	instance, err := multipass.GetInstance(name)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	slog.Debug("get instance", slog.Any("instance", instance))
	if instance == nil {
		return fmt.Errorf("instance not found: %s", name)
	}

	err = multipass.Exec(name, "/opt/csr.sh")
	if err != nil {
		return fmt.Errorf("failed to execute csr.sh: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "kom")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	err = multipass.Transfer(name, "/home/ubuntu/.kube/config", tempDir)
	if err != nil {
		return fmt.Errorf("failed to transfer kubeconfig file: %w", err)
	}

	kubeDir := os.Getenv("HOME") + "/.kube"
	err = os.MkdirAll(kubeDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .kube directory: %w", err)
	}

	mergedConfig, err := kubernetes.MergeKubeconfig([]string{kubeDir + "/config", tempDir + "/config"})
	if err != nil {
		return fmt.Errorf("failed to merge kubeconfig files: %w", err)
	}

	mergedConfig.Clusters["kubernetes"].Server = fmt.Sprintf("https://%s:6443", instance.Ipv4[0])

	json, err := runtime.Encode(clientcmdlatest.Codec, mergedConfig)
	if err != nil {
		return fmt.Errorf("failed to encode kubeconfig to JSON: %w", err)
	}
	output, err := yaml.JSONToYAML(json)
	if err != nil {
		return fmt.Errorf("failed to convert JSON to YAML: %w", err)
	}

	err = os.WriteFile(kubeDir+"/config", output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig file: %w", err)
	}

	return nil
}
