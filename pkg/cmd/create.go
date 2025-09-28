package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/provisioner"
)

const (
	defaultInstanceCpus   = "2"
	defaultInstanceMemory = "4G"
	defaultInstanceDisk   = "10G"
	defaultK8sVersion     = "v1.34.0"
	defaultInstanceImage  = "24.04"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster, master or worker node",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a new cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.CreateCluster(cmd.Flag("cluster-name").Value.String(),
			provisioner.ClusterConfig{},
			provisioner.MasterConfig{
				InstanceConfig: provisioner.InstanceConfig{
					Name:       "master",
					CPUs:       defaultInstanceCpus,
					Memory:     defaultInstanceMemory,
					Disk:       defaultInstanceDisk,
					K8sVersion: defaultK8sVersion,
					Image:      defaultInstanceImage,
				},
				IsRegisterNode: true,
			},
			provisioner.WorkerConfig{
				InstanceConfig: provisioner.InstanceConfig{
					Name:       "worker",
					CPUs:       defaultInstanceCpus,
					Memory:     defaultInstanceMemory,
					Disk:       defaultInstanceDisk,
					K8sVersion: defaultK8sVersion,
					Image:      defaultInstanceImage,
				},
				IsJoinCluster: true,
			},
		)
	},
}

var createWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Create a new worker node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provisioner.CreateWorker(cmd.Flag("cluster-name").Value.String(), getWorkerConfig(cmd))
	},
}

func getWorkerConfig(cmd *cobra.Command) provisioner.WorkerConfig {
	join, _ := cmd.Flags().GetBool("join")

	return provisioner.WorkerConfig{
		InstanceConfig: provisioner.InstanceConfig{
			Name:       cmd.Flag("name").Value.String(),
			CPUs:       cmd.Flag("cpus").Value.String(),
			Memory:     cmd.Flag("memory").Value.String(),
			Disk:       cmd.Flag("disk").Value.String(),
			K8sVersion: cmd.Flag("k8s-version").Value.String(),
			Image:      cmd.Flag("image").Value.String(),
		},
		IsJoinCluster: join,
	}
}

func defineWorkerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "Name of the instance. If not provided, a random name will be generated")
	cmd.Flags().StringP("cpus", "c", defaultInstanceCpus, "Number of CPUs")
	cmd.Flags().StringP("memory", "m", defaultInstanceMemory, "Amount of memory")
	cmd.Flags().StringP("disk", "d", defaultInstanceDisk, "Amount of disk space")
	cmd.Flags().StringP("k8s-version", "k", defaultK8sVersion, "Kubernetes version")
	cmd.Flags().StringP("image", "i", defaultInstanceImage, "Image to use for the VM")
	cmd.Flags().BoolP("join", "j", true, "Join the cluster")
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createClusterCmd)
	createCmd.AddCommand(createWorkerCmd)

	createCmd.PersistentFlags().StringP("cluster-name", "", "kubernetes", "Name of the cluster")
	defineWorkerFlags(createWorkerCmd)
}
