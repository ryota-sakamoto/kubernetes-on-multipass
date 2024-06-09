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

var createMasterCmd = &cobra.Command{
	Use:   "master",
	Short: "Create a new master",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.CreateMaster(getProvisionerConfig(cmd))
	},
}

var createWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Create a new worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.CreateWorker(getProvisionerConfig(cmd))
	},
}

func getProvisionerConfig(cmd *cobra.Command) provisioner.Config {
	return provisioner.Config{
		Name:       cmd.Flag("prefix").Value.String() + cmd.Flag("name").Value.String(),
		CPUs:       cmd.Flag("cpus").Value.String(),
		Memory:     cmd.Flag("memory").Value.String(),
		Disk:       cmd.Flag("disk").Value.String(),
		K8sVersion: cmd.Flag("k8s-version").Value.String(),
		Image:      cmd.Flag("image").Value.String(),
	}
}

func defineCommonFlags(cmd *cobra.Command, name string) {
	cmd.Flags().StringP("prefix", "p", "kubernetes-", "Prefix for the instance")
	cmd.Flags().StringP("name", "n", name, "Name of the instance")
	cmd.Flags().StringP("cpus", "c", "2", "Number of CPUs")
	cmd.Flags().StringP("memory", "m", "4G", "Amount of memory")
	cmd.Flags().StringP("disk", "d", "10G", "Amount of disk space")
	cmd.Flags().StringP("k8s-version", "k", "v1.30.0", "Kubernetes version")
	cmd.Flags().StringP("image", "i", "22.04", "Image to use for the VM")

}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createMasterCmd)
	createCmd.AddCommand(createWorkerCmd)
	defineCommonFlags(createMasterCmd, "master")
	defineCommonFlags(createWorkerCmd, "worker")
}
