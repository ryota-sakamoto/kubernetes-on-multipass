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

```bash
$ kubectl get node -o wide
NAME                STATUS   ROLES           AGE   VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION       CONTAINER-RUNTIME
kubernetes-master   Ready    control-plane   29m   v1.30.0   10.81.43.68    <none>        Ubuntu 22.04.4 LTS   5.15.0-105-generic   containerd://1.7.0
kubernetes-worker   Ready    <none>          16m   v1.30.0   10.81.43.118   <none>        Ubuntu 22.04.4 LTS   5.15.0-105-generic   containerd://1.7.0
```

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

| Env Variable | Description | Default Value |
| - | - | - |
| MASTER_INSTANCE | Name of the master instance | master |
| WORKER_INSTANCE | Name of the worker instance | worker |
| INSTANCE_NAME_PREFIX | Prefix to be used for the instance names | kubernetes- |
| CPU | CPU count of the instance | 2 |
| MEMORY | Memory of the instance | 2G |
| DISK | Disk size of the instance | 10G |
| IMAGE | Image of the instance | 22.04 |
| KUBERNETES_VERSION | Version of Kubernetes to install | v1.30.0 |

These variables can be overridden by setting them in the shell or by editing the Makefile directly.

To create a cluster with Kubernetes version v1.27.1 instead of the default version specified in the Makefile, you can run the following command:

```bash
make create-cluster KUBERNETES_VERSION=v1.27.1
```

To create a worker instance with a custom name and join the cluster, you can run the following command:

```bash
make create-worker WORKER_INSTANCE=worker2
make join-worker WORKER_INSTANCE=worker2
```
