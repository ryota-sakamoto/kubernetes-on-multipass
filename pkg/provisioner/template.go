package provisioner

import "github.com/ryota-sakamoto/kubernetes-on-multipass/pkg/cloudinit"

var baseTemplate = cloudinit.Config{
	Packages: []string{
		"apt-transport-https",
		"ca-certificates",
		"curl",
		"gpg",
		"jq",
		"conntrack",
	},
	WriteFiles: []cloudinit.WriteFile{
		{
			Content: `overlay
br_netfilter`,
			Path: "/etc/modules-load.d/k8s.conf",
		},
		{
			Content: `net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1`,
			Path: "/etc/sysctl.d/k8s.conf",
		},
		{
			Content: `#!/bin/bash
cd /usr/bin

CRICTL_VERSION="v1.28.0"
curl -L "https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-{{ .Arch }}.tar.gz" | sudo tar -C /usr/bin -xz

sudo curl -L --remote-name-all https://dl.k8s.io/release/{{ .KubernetesVersion }}/bin/linux/{{ .Arch }}/{kubeadm,kubelet,kubectl} -O
sudo chmod +x {kubeadm,kubelet,kubectl}

RELEASE_VERSION="v0.16.2"
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${RELEASE_VERSION}/cmd/krel/templates/latest/kubelet/kubelet.service" | sudo tee /usr/lib/systemd/system/kubelet.service
sudo mkdir -p /usr/lib/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${RELEASE_VERSION}/cmd/krel/templates/latest/kubeadm/10-kubeadm.conf" | sudo tee /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf`,
			Path:        "/opt/tools.sh",
			Permissions: "0755",
		},
		{
			Content: `#!/bin/sh
wget https://github.com/containerd/containerd/releases/download/v1.7.0/containerd-1.7.0-linux-{{ .Arch }}.tar.gz
tar Cxzvf /usr/local containerd-1.7.0-linux-{{ .Arch }}.tar.gz

wget https://raw.githubusercontent.com/containerd/containerd/v1.7.0/containerd.service
mv containerd.service /etc/systemd/system/containerd.service
systemctl daemon-reload
systemctl enable --now containerd`,
			Path:        "/opt/containerd.sh",
			Permissions: "0755",
		},
		{
			Content: `version = 2
[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
    [plugins."io.containerd.grpc.v1.cri".containerd]
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          runtime_type = "io.containerd.runc.v2"
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
            SystemdCgroup = true`,
			Path: "/etc/containerd/config.toml",
		},
		{
			Content: `#!/bin/sh
wget https://github.com/opencontainers/runc/releases/download/v1.1.4/runc.{{ .Arch }}
install -m 755 runc.{{ .Arch }} /usr/local/sbin/runc`,
			Path:        "/opt/runc.sh",
			Permissions: "0755",
		},
		{
			Content: `#!/bin/sh
wget https://github.com/containernetworking/plugins/releases/download/v1.2.0/cni-plugins-linux-{{ .Arch }}-v1.2.0.tgz
mkdir -p /opt/cni/bin
tar Cxzvf /opt/cni/bin cni-plugins-linux-{{ .Arch }}-v1.2.0.tgz`,
			Path:        "/opt/cni.sh",
			Permissions: "0755",
		},
		{
			Content: "KUBELET_EXTRA_ARGS=--cgroup-driver=systemd",
			Path:    "/etc/default/kubelet",
		},
	},
	RunCmds: []string{
		"modprobe -a overlay br_netfilter",
		"sysctl --system",
		"/opt/tools.sh",
		"/opt/containerd.sh",
		"/opt/runc.sh",
		"/opt/cni.sh",
	},
}

var masterTemplate = cloudinit.Config{
	WriteFiles: []cloudinit.WriteFile{
		{
			Content: `apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
nodeRegistration:
  ignorePreflightErrors:
    - NumCPU
    - Mem
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
registerNode: {{ .RegisterNode }}
---
apiVersion: kubeadm.k8s.io/v1beta3
kind: ClusterConfiguration
kubernetesVersion: {{ .KubernetesVersion }}
networking:
  serviceSubnet: "172.20.0.0/16"
  podSubnet: "192.168.0.0/16"
  dnsDomain: "cluster.local"
apiServer:
  extraArgs:
    disable-admission-plugins: CertificateSubjectRestriction`,
			Path: "/opt/kubeadm/init.yaml",
		},
		{
			Content: `#!/bin/sh
kubeadm init --config /opt/kubeadm/init.yaml`,

			Path:        "/opt/kubeadm.sh",
			Permissions: "0755",
		},
		{
			Content: `#!/bin/sh
openssl genrsa -out admin.pem 2048
openssl req -new -key admin.pem -out admin.csr -subj "/O=system:masters/CN=admin"
CSR=$(cat admin.csr | base64 | tr -d "\n")
sed "s/__CSR__/${CSR}/" /opt/csr.yaml > /tmp/csr.yaml

sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl delete -f /tmp/csr.yaml
sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl apply -f /tmp/csr.yaml
sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl certificate approve admin
sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl get csr admin -o json | jq -r .status.certificate | base64 --decode > admin.crt

kubectl config set-cluster kubernetes --certificate-authority=/etc/kubernetes/pki/ca.crt --embed-certs --server=https://localhost:6443
kubectl config set-credentials admin --client-certificate=admin.crt --client-key=admin.pem --embed-certs=true
kubectl config set-context kubernetes --cluster=kubernetes --user=admin
kubectl config use-context kubernetes`,
			Path:        "/opt/csr.sh",
			Permissions: "0755",
		},
		{
			Content: `apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: admin
spec:
  signerName: kubernetes.io/kube-apiserver-client
  groups:
  - system:masters
  request: __CSR__
  usages:
  - digital signature
  - key encipherment
  - client auth`,
			Path: "/opt/csr.yaml",
		},
	},
	RunCmds: []string{
		"/opt/kubeadm.sh",
	},
}
