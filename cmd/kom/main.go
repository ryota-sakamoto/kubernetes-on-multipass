package main

import (
	"log/slog"
	"os"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cmd"
)

func init() {
	logLevel := slog.LevelInfo
	if l, b := os.LookupEnv("LOG_LEVEL"); b {
		logLevel.UnmarshalText([]byte(l))
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))
}

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
