package main

import (
	"log/slog"
	"os"

	"github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
