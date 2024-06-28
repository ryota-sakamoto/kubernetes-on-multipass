package multipass

import (
	"encoding/json"
	"log/slog"
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

func commandWrapper(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	slog.Debug("execute command", slog.String("command", cmd.String()))

	return cmd
}

func ListInstances() (*Instances, error) {
	cmd := commandWrapper("multipass", "list", "--format", "json")

	output, err := cmd.Output()
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
	cmd := commandWrapper("multipass", "launch", config.Image, "--name", config.Name, "-c", config.CPUs, "-m", config.Memory, "-d", config.Disk, "--cloud-init", "-")
	cmd.Stdin = strings.NewReader(cloudinit)

	_, err := cmd.Output()
	return err
}

func DeleteInstance(name string) error {
	_, err := commandWrapper("multipass", "delete", name).Output()
	if err != nil {
		return err
	}

	return nil
}

func Purge() error {
	_, err := commandWrapper("multipass", "purge").Output()
	if err != nil {
		return err
	}

	return nil
}

func Exec(name string, command string) (string, error) {
	args := []string{"exec", name, "--"}
	args = append(args, strings.Fields(command)...)

	cmd := commandWrapper("multipass", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Transfer(name string, from string, to string) error {
	args := []string{"transfer", name + ":" + from, to}

	cmd := commandWrapper("multipass", args...)
	_, err := cmd.Output()
	return err
}
