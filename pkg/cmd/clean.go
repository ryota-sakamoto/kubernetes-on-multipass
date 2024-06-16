package cmd

import (
	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/provisioner"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up a cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.Clean(cmd.Flag("cluster-name").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("cluster-name", "", "kubernetes", "Cluster name")
}
