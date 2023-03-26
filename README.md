Kubernetes on Multipass
===

This repository contains a Makefile that automates the deployment of a Kubernetes cluster on [Multipass](https://multipass.run/), a lightweight VM manager for Linux, macOS, and Windows.

## Prerequisites

Before using this Makefile, you will need to have the following installed on your machine:

- Multipass
- kubectl
- jq, a command-line JSON processor
- Helm

## Usage

To create a Kubernetes cluster on Multipass and install CNI with Helm, simply run the following command:

```bash
make create-cluster
```

This will launch two Multipass instances, master and worker, and generate a Kubernetes configuration file. It will then join the worker to the cluster and install the Cilium CNI using Helm.

You can also run each target separately using the following commands:

```bash
# Launch the master and worker instances
make create-master
make create-worker

# Join the worker to the cluster
make join-worker

# Generate a Kubernetes configuration file
make generate-kubeconfig

# Install CNI with Helm
make install-cni
```

To open a shell session on the master instance, run:

```bash
make shell-master
```

To clean up the resources, run:

```bash
make clean
```

## Customization

The Makefile provides several variables that can be customized:

- `MASTER_INSTANCE` and `WORKER_INSTANCE`: Names to be given to the master and worker instances.
- `KUBERNETES_VERSION`: Version of Kubernetes to install on the instances. This variable is used to set the version of the `kubelet`, `kubeadm`, and `kubectl` packages.

These variables can be overridden by setting them in the shell or by editing the Makefile directly.

For example, to create a cluster with Kubernetes version 1.25.8-00 instead of the default version specified in the Makefile, you can run the following command:

```bash
make create-cluster KUBERNETES_VERSION=1.25.8-00
```
