package cloudinit

import (
	"bytes"
	"text/template"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Packages   []string    `yaml:"packages"`
	WriteFiles []WriteFile `yaml:"write_files"`
	RunCmds    []string    `yaml:"runcmd"`

	Vars map[string]any `yaml:"-"`
}

type WriteFile struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Permissions string `yaml:"permissions,omitempty"`
}

func (c *Config) Generate() (string, error) {
	buff := &bytes.Buffer{}
	if err := yaml.NewEncoder(buff).Encode(c); err != nil {
		return "", err
	}

	temp, err := template.New("cloud-config").Parse(`#cloud-config
package_update: true

` + buff.String())
	if err != nil {
		return "", err
	}

	output := &bytes.Buffer{}
	if err := temp.Execute(output, c.Vars); err != nil {
		return "", err
	}

	return output.String(), nil
}
