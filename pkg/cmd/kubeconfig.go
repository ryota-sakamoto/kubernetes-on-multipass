package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/provisioner"
)

var generateKubeconfigCmd = &cobra.Command{
	Use:   "generate-kubeconfig",
	Short: "Generate a kubeconfig file for the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.GenerateKubeconfig("kubernetes-master")
	},
}

func init() {
	rootCmd.AddCommand(generateKubeconfigCmd)
}
