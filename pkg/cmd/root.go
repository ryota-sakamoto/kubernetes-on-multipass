package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "kom",
	Short:         "kom is a CLI tool to deploy Kubernetes on Multipass",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel := slog.LevelInfo
		if b, _ := cmd.Flags().GetBool("debug"); b {
			logLevel = slog.LevelDebug
		}

		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})))
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug mode")
}
