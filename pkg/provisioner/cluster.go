package provisioner

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/kubernetes"
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/multipass"
)

type ClusterConfig struct {
	Name string
}

type InstanceConfig struct {
	Name          string
	CPUs          string
	Memory        string
	Disk          string
	K8sVersion    string
	Image         string
	IsJoinCluster bool
}

func CreateCluster(clusterName string, config ClusterConfig) error {
	slog.Debug("create cluster", slog.String("clusterName", clusterName), slog.Any("config", config))

	return nil
}

func CreateMaster(clusterName string, config InstanceConfig) error {
	slog.Debug("create master", slog.String("clusterName", clusterName), slog.Any("config", config))

	config.Name = "master"
	_, err := LaunchInstance(clusterName, config, GetMasterTemplate())
	if err != nil {
		return fmt.Errorf("failed to launch instance: %w", err)
	}

	return nil
}

func CreateWorker(clusterName string, config InstanceConfig) error {
	slog.Debug("create worker", slog.String("clusterName", clusterName), slog.Any("config", config))

	instanceName, err := LaunchInstance(clusterName, config, GetWorkerTemplate())
	if err != nil {
		return fmt.Errorf("failed to launch instance: %w", err)
	}

	if config.IsJoinCluster {
		return JoinCluster(clusterName, instanceName)
	}

	return nil
}

func JoinCluster(clusterName, name string) error {
	slog.Debug("join cluster", slog.String("clusterName", clusterName), slog.String("name", name))

	masterName := clusterName + "-master"
	joinCommand, err := multipass.Exec(masterName, "sudo kubeadm token create --print-join-command")
	if err != nil {
		return fmt.Errorf("failed to get join command: %w", err)
	}

	_, err = multipass.Exec(name, fmt.Sprintf("sudo %s", joinCommand))
	if err != nil {
		return fmt.Errorf("failed to join cluster: %w", err)
	}

	return nil
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

	_, err = multipass.Exec(name, "/opt/csr.sh")
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
