package cloudinit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"
)

func TestGenerate(t *testing.T) {
	result, err := (&cloudinit.Config{
		Packages: []string{
			"git",
			"vim",
		},
		WriteFiles: []cloudinit.WriteFile{
			{
				Content: `file {{ .Content }}`,
				Path:    "/opt/ok",
			},
			{
				Content:     `echo ok`,
				Path:        "/use/bin/ok",
				Permissions: "0755",
			},
		},
		RunCmds: []string{
			"/use/bin/ok",
		},
		Vars: map[string]any{
			"Content": "content",
		},
	}).Generate()
	assert.NoError(t, err)

	expectedResult := `#cloud-config
package_update: true

packages:
- git
- vim
write_files:
- path: /opt/ok
  content: file content
- path: /use/bin/ok
  content: echo ok
  permissions: "0755"
runcmd:
- /use/bin/ok
`
	assert.Equal(t, expectedResult, result)
}
