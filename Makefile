INSTANCE_NAME_PREFIX := kubernetes-
MASTER_INSTANCE := master
WORKER_INSTANCE := worker

MASTER_INSTANCE_NAME := $(INSTANCE_NAME_PREFIX)$(MASTER_INSTANCE)
WORKER_INSTANCE_NAME := $(INSTANCE_NAME_PREFIX)$(WORKER_INSTANCE)

KUBERNETES_VERSION := v1.30.0

error:
	exit 1

create-cluster: create-master create-worker generate-kubeconfig join-worker install-cni

create-master:
	cat cloud-init-master.yaml | sed "s/__KUBERNETES_VERSION__/$(KUBERNETES_VERSION)/" | multipass launch 22.04 --name $(MASTER_INSTANCE_NAME) -c 2 -m 1G -d 10G --cloud-init -

create-worker:
	cat cloud-init-worker.yaml | sed "s/__KUBERNETES_VERSION__/$(KUBERNETES_VERSION)/" | multipass launch 22.04 --name $(WORKER_INSTANCE_NAME) -c 2 -m 1G -d 10G --cloud-init -

join-worker:
	$(eval JOIN_COMMAND := $(shell multipass exec $(MASTER_INSTANCE_NAME) -- sudo kubeadm token create --print-join-command))
	multipass exec $(WORKER_INSTANCE_NAME) -- sudo $(JOIN_COMMAND)

shell-master:
	multipass shell $(MASTER_INSTANCE_NAME)

generate-kubeconfig:
	multipass exec $(MASTER_INSTANCE_NAME) -- /opt/csr.sh
	multipass transfer $(MASTER_INSTANCE_NAME):/home/ubuntu/.kube/config .
	KUBECONFIG=config:~/.kube/config kubectl config view --flatten > ~/.kube/config
	rm config

	$(eval IP := $(shell multipass info $(MASTER_INSTANCE_NAME) --format json | jq .info | jq -r '.["$(MASTER_INSTANCE_NAME)"].ipv4[0]'))
	kubectl config set-cluster kubernetes --server=https://$(IP):6443

install-cni:
	helm repo add cilium https://helm.cilium.io/
	helm install cilium cilium/cilium --version 1.13.1 --namespace kube-system --set ipam.mode=kubernetes

clean:
	-multipass list --format json | jq -r .list[].name | grep "$(INSTANCE_NAME_PREFIX)" | xargs multipass delete
	-multipass purge

update-kubernetes-version:
	$(eval LATEST_VRESION := $(shell curl -s -H "Accept: application/vnd.github+json" https://api.github.com/repos/kubernetes/kubernetes/releases/latest | jq -r .tag_name | sed "s/v//"))
	sed -i "s/${KUBERNETES_VERSION}/${LATEST_VRESION}-00/" Makefile README.md
