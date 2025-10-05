package cmd

import (
	"log/slog"
	"os"
	"strings"

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

		logFormat, _ := cmd.Flags().GetString("log-format")
		opts := &slog.HandlerOptions{Level: logLevel}

		if strings.ToLower(logFormat) == "json" {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug logging")
	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json)")
}
