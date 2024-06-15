package kubernetes_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/kubernetes"
)

const baseKubeconfig = `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: aG9nZQ==
    server: https://192.168.205.%[1]d:6443
  name: kubernetes%[1]d
contexts:
- context:
    cluster: kubernetes%[1]d
    user: admin%[1]d
  name: kubernetes%[1]d
current-context: kubernetes%[1]d
users:
- name: admin%[1]d
  user:
    client-certificate-data: dGVzdAo=
    client-key-data: YWJjCg==
`

func TestMergeKubeconfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kom")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	makeTempKubeconfig(t, tempDir+"/kubeconfig1", fmt.Sprintf(baseKubeconfig, 1))
	makeTempKubeconfig(t, tempDir+"/kubeconfig2", fmt.Sprintf(baseKubeconfig, 2))

	mergedConfig, err := kubernetes.MergeKubeconfig([]string{tempDir + "/kubeconfig1", tempDir + "/kubeconfig2"})
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, &api.Config{
		APIVersion: "",
		Kind:       "",
		Preferences: api.Preferences{
			Colors:     false,
			Extensions: make(map[string]runtime.Object),
		},
		Clusters: map[string]*api.Cluster{
			"kubernetes1": {
				LocationOfOrigin:         tempDir + "/kubeconfig1",
				Server:                   "https://192.168.205.1:6443",
				CertificateAuthorityData: []byte("hoge"),
				Extensions:               make(map[string]runtime.Object),
			},
			"kubernetes2": {
				LocationOfOrigin:         tempDir + "/kubeconfig2",
				Server:                   "https://192.168.205.2:6443",
				CertificateAuthorityData: []byte("hoge"),
				Extensions:               make(map[string]runtime.Object),
			},
		},
		Contexts: map[string]*api.Context{
			"kubernetes1": {
				LocationOfOrigin: tempDir + "/kubeconfig1",
				Cluster:          "kubernetes1",
				AuthInfo:         "admin1",
				Extensions:       make(map[string]runtime.Object),
			},
			"kubernetes2": {
				LocationOfOrigin: tempDir + "/kubeconfig2",
				Cluster:          "kubernetes2",
				AuthInfo:         "admin2",
				Extensions:       make(map[string]runtime.Object),
			},
		},
		CurrentContext: "kubernetes1",
		AuthInfos: map[string]*api.AuthInfo{
			"admin1": {
				LocationOfOrigin:      tempDir + "/kubeconfig1",
				ClientCertificateData: []byte("test\n"),
				ClientKeyData:         []byte("abc\n"),
				Extensions:            make(map[string]runtime.Object),
			},
			"admin2": {
				LocationOfOrigin:      tempDir + "/kubeconfig2",
				ClientCertificateData: []byte("test\n"),
				ClientKeyData:         []byte("abc\n"),
				Extensions:            make(map[string]runtime.Object),
			},
		},
		Extensions: make(map[string]runtime.Object),
	}, mergedConfig)
}

func makeTempKubeconfig(t *testing.T, path, content string) {
	t.Helper()

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
