Kubernetes on Multipass
===

This repository contains a CLI that automates the deployment of a Kubernetes cluster on [Multipass](https://multipass.run/), a lightweight VM manager for Linux, macOS, and Windows.

## Prerequisites

Before using this tool, you will need to have the following installed on your machine:

- Multipass

## Installation

```bash
go install github.com/ryota-sakamoto/kubernetes-on-multipass/cmd/kom@latest
```

## Usage

To create a Kubernetes cluster on Multipass and install CNI with Helm, simply run the following command:

```bash
kom create cluster
```

This will launch two Multipass instances, a master and a worker, and generate a Kubernetes configuration file. It will then join the worker to the cluster and install the Cilium CNI using Helm.

```bash
$ kubectl get node -o wide
NAME                STATUS   ROLES           AGE     VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
kubernetes-master   Ready    control-plane   2m49s   v1.34.0   192.168.205.73   <none>        Ubuntu 24.04.3 LTS   6.8.0-71-generic   containerd://1.7.0
kubernetes-worker   Ready    <none>          72s     v1.34.0   192.168.205.74   <none>        Ubuntu 24.04.3 LTS   6.8.0-71-generic   containerd://1.7.0
```

You can also run each command separately:

```bash
# Launch the worker instances
kom create worker

# Generate a Kubernetes configuration file
kom generate kubeconfig

# Install CNI with Helm
kom install cni
```

To run a command on the instance, run:

```bash
kom exec master -- bash
```

To clean up the resources, run:

```bash
kom clean
```

## Shell Completion

kom supports generating shell completion scripts for various shells via cobra.

### Bash

To load completions for the current session:

```bash
source <(kom completion bash)
```

To make completions available for all sessions, you can add the completion script to your system's bash completion directory:

```bash
kom completion bash > /etc/bash_completion.d/kom # Linux
kom completion bash > /usr/local/etc/bash_completion.d/kom # macOS
```

### Fish

To load completions for the current session:

```bash
kom completion fish | source
```

To make completions available for all sessions, save the completion script to your fish completions directory:

```bash
kom completion fish > ~/.config/fish/completions/kom.fish
```
