package cloudinit

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Packages   []string    `yaml:"packages"`
	WriteFiles []WriteFile `yaml:"write_files"`
	RunCmds    []string    `yaml:"runcmd"`
}

type WriteFile struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Permissions string `yaml:"permissions,omitempty"`
}

func Generate(c Config, vars map[string]string) (string, error) {
	buff := &bytes.Buffer{}
	if err := yaml.NewEncoder(buff).Encode(c); err != nil {
		return "", err
	}

	temp, err := template.New("cloud-config").Parse(buff.String())
	if err != nil {
		return "", err
	}

	if err := temp.Execute(buff, vars); err != nil {
		return "", err
	}

	return fmt.Sprintf(`#cloud-config
package_update: true

%s`, buff.String()), nil
}
