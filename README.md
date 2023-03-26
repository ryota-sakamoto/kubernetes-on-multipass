kubernetes-on-multipass
===

## master

```
$ make master
$ make kubeconfig
```

## worker

```
$ make worker
$ make join
```

## CNI

```
$ helm repo add cilium https://helm.cilium.io/
$ helm install cilium cilium/cilium --version 1.13.1 --namespace kube-system --set ipam.mode=kubernetes
```
