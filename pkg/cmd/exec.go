package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/multipass"
)

var execCmd = &cobra.Command{
	Use:   "exec <node-name> -- [command]",
	Short: "Execute a command in an instance",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeName := args[0]
		commandSeparator := cmd.ArgsLenAtDash()
		if commandSeparator == -1 || commandSeparator == 0 {
			return errors.New("a command is required after '--'")
		}
		command := strings.Join(args[commandSeparator:], " ")

		clusterName, err := cmd.Flags().GetString("cluster-name")
		if err != nil {
			return err
		}

		instanceName := fmt.Sprintf("%s-%s", clusterName, nodeName)
		slog.Info("Executing command in instance", slog.String("instance", instanceName), slog.String("command", command))
		return multipass.ExecInteractive(instanceName, command)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("cluster-name", "", "kubernetes", "Name of the cluster")
}
