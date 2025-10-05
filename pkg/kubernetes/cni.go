package kubernetes

import (
	"fmt"
	"log/slog"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

const (
	ciliumRepoName    = "cilium"
	ciliumRepoURL     = "https://helm.cilium.io/"
	ciliumChartName   = "cilium"
	ciliumReleaseName = "cilium"
	ciliumNamespace   = "kube-system"
)

type CNI struct {
	actionConfig *action.Configuration
	settings     *cli.EnvSettings
}

func NewCNI(kubeconfigPath, context string) (*CNI, error) {
	settings := cli.New()
	settings.KubeConfig = kubeconfigPath
	settings.KubeContext = context

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), slog.Debug); err != nil {
		return nil, fmt.Errorf("failed to init action config: %w", err)
	}

	return &CNI{
		actionConfig: actionConfig,
		settings:     settings,
	}, nil
}

func (c *CNI) InstallCilium() error {
	slog.Info("Installing Cilium CNI via Helm...")

	install := action.NewInstall(c.actionConfig)
	install.ReleaseName = ciliumReleaseName
	install.Namespace = ciliumNamespace
	install.RepoURL = ciliumRepoURL

	chartPath, err := install.LocateChart(ciliumChartName, c.settings)
	if err != nil {
		return fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	if _, err := install.Run(chart, nil); err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}

	slog.Info("Cilium installed successfully")
	return nil
}
