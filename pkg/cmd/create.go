package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/provisioner"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster, master or worker node",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Create a new worker node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.CreateWorker(cmd.Flag("cluster-name").Value.String(), getProvisionerConfig(cmd))
	},
}

func getProvisionerConfig(cmd *cobra.Command) provisioner.Config {
	return provisioner.Config{
		Name:       cmd.Flag("name").Value.String(),
		CPUs:       cmd.Flag("cpus").Value.String(),
		Memory:     cmd.Flag("memory").Value.String(),
		Disk:       cmd.Flag("disk").Value.String(),
		K8sVersion: cmd.Flag("k8s-version").Value.String(),
		Image:      cmd.Flag("image").Value.String(),
	}
}

func defineCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("cluster-name", "", "kubernetes", "Name of the cluster")

	cmd.Flags().StringP("name", "n", "", "Name of the instance. If not provided, a random name will be generated")
	cmd.Flags().StringP("cpus", "c", "2", "Number of CPUs")
	cmd.Flags().StringP("memory", "m", "4G", "Amount of memory")
	cmd.Flags().StringP("disk", "d", "10G", "Amount of disk space")
	cmd.Flags().StringP("k8s-version", "k", "v1.30.0", "Kubernetes version")
	cmd.Flags().StringP("image", "i", "22.04", "Image to use for the VM")
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createWorkerCmd)
	defineCommonFlags(createWorkerCmd)
}
