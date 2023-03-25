error:
	exit 1

master:
	multipass launch 22.04 --name master --cloud-init cloud-init.yaml

shell:
	multipass shell master

clean:
	multipass delete master
	multipass purge
