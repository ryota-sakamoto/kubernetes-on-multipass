package main

import "github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
