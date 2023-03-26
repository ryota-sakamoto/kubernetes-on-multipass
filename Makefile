MASTER_INSTANCE := master
WORKER_INSTANCE := worker

KUBERNETES_VERSION := 1.26.3-00

error:
	exit 1

create-cluster: create-master create-worker generate-kubeconfig join-worker install-cni

create-master:
	cat cloud-init-master.yaml | sed "s/__KUBERNETES_VERSION__/$(KUBERNETES_VERSION)/" | multipass launch 22.04 --name $(MASTER_INSTANCE) -c 2 -m 1G -d 10G --cloud-init -

create-worker:
	cat cloud-init-worker.yaml | sed "s/__KUBERNETES_VERSION__/$(KUBERNETES_VERSION)/" | multipass launch 22.04 --name $(WORKER_INSTANCE) -c 2 -m 1G -d 10G --cloud-init -

join-worker:
	$(eval JOIN_COMMAND := $(shell multipass exec $(MASTER_INSTANCE) -- sudo kubeadm token create --print-join-command))
	multipass exec $(WORKER_INSTANCE) -- sudo $(JOIN_COMMAND)

shell-master:
	multipass shell $(MASTER_INSTANCE)

generate-kubeconfig:
	multipass exec $(MASTER_INSTANCE) -- /opt/csr.sh
	multipass transfer $(MASTER_INSTANCE):/home/ubuntu/.kube/config .
	KUBECONFIG=config:~/.kube/config kubectl config view --flatten > ~/.kube/config
	rm config

	$(eval IP := $(shell multipass info $(MASTER_INSTANCE) --format json | jq -r .info.master.ipv4[0]))
	kubectl config set-cluster kubernetes --server=https://$(IP):6443

install-cni:
	helm repo add cilium https://helm.cilium.io/
	helm install cilium cilium/cilium --version 1.13.1 --namespace kube-system --set ipam.mode=kubernetes

clean:
	-multipass delete $(WORKER_INSTANCE)
	-multipass delete $(MASTER_INSTANCE)
	-multipass purge
