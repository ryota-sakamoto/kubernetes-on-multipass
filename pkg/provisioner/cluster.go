package provisioner

import (
	"fmt"
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

func GenerateKubeconfig(name string) error {
	instance, err := multipass.GetInstance(name)
	if err != nil {
		return err
	}

	err = multipass.Exec(name, "/opt/csr.sh")
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "kom")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = multipass.Transfer(name, "/home/ubuntu/.kube/config", tempDir)
	if err != nil {
		return err
	}

	kubeDir := os.Getenv("HOME") + "/.kube"
	err = os.MkdirAll(kubeDir, 0755)
	if err != nil {
		return err
	}

	mergedConfig, err := kubernetes.MergeKubeconfig([]string{kubeDir + "/config", tempDir + "/config"})
	if err != nil {
		return err
	}

	mergedConfig.Clusters["kubernetes"].Server = fmt.Sprintf("https://%s:6443", instance.Ipv4[0])

	json, err := runtime.Encode(clientcmdlatest.Codec, mergedConfig)
	if err != nil {
		return nil
	}
	output, err := yaml.JSONToYAML(json)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}

	err = os.WriteFile(kubeDir+"/config", output, 0644)
	if err != nil {
		return err
	}

	return nil
}
