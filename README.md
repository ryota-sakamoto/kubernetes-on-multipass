kubernetes-on-multipass
===

## setup

```
$ make master
$ make kubeconfig
```

## CNI

```
$ helm install cilium cilium/cilium --version 1.13.1 --namespace kube-system --set ipam.mode=kubernetes
```
