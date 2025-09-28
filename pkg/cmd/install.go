package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/provisioner"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install components to the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cniCmd = &cobra.Command{
	Use:   "cni",
	Short: "Install CNI to the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.InstallCNI(cmd.Flag("cluster-name").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(cniCmd)
	installCmd.PersistentFlags().StringP("cluster-name", "", "kubernetes", "Name of the cluster")
}
