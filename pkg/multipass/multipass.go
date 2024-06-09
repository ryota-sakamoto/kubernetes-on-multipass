package multipass

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type Instances struct {
	List []Instance `json:"list"`
}

type Instance struct {
	Name    string   `json:"name"`
	State   string   `json:"state"`
	Ipv4    []string `json:"ipv4"`
	Release string   `json:"release"`
}

func ListInstances() (*Instances, error) {
	output, err := exec.Command("multipass", "list", "--format", "json").Output()
	if err != nil {
		return nil, err
	}

	var instances *Instances
	if err := json.Unmarshal(output, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

func GetInstance(name string) (*Instance, error) {
	instances, err := ListInstances()
	if err != nil {
		return nil, err
	}

	for _, instance := range instances.List {
		if instance.Name == name {
			return &instance, nil
		}
	}

	return nil, nil
}

type InstanceConfig struct {
	Name   string
	CPUs   string
	Memory string
	Disk   string
	Image  string
}

func LaunchInstance(config InstanceConfig, cloudinit string) error {
	cmd := exec.Command("multipass", "launch", config.Image, "--name", config.Name, "-c", config.CPUs, "-m", config.Memory, "-d", config.Disk, "--cloud-init", "-")
	cmd.Stdin = strings.NewReader(cloudinit)

	_, err := cmd.Output()
	return err
}
