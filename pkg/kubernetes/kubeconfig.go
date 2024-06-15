package kubernetes

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func MergeKubeconfig(path []string) (*api.Config, error) {
	loadingRules := clientcmd.ClientConfigLoadingRules{
		Precedence: path,
	}

	mergedConfig, err := loadingRules.Load()
	if err != nil {
		return nil, err
	}

	return mergedConfig, nil
}
