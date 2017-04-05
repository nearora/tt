#!/bin/bash

usage() {
	echo "Usage: $0 start|stop|status" >&2
}

containerIds() {
	docker ps -aq 2>/dev/null
	if [ $? -ne 0 ] ; then
		echo "Error communicating with docker." >&2
	fi
}

numContainers() {
	containerIds | wc -l
}

if [ $# -lt 1 ] ; then
	usage
	exit 1
fi

if [ "$1x" = "statusx" ] ; then
	n=$(numContainers)
	echo "$n containers running..."
elif [ "$1x" = "stopx" ] ; then
	if [ $(numContainers) -ge 1 ] ; then
		echo -n "Killing and removing all docker containers on the system... "
		docker rm -f $(containerIds) >/dev/null 2>&1
		if [ $? -ne 0 ] ; then
			echo "failed!"
		else
			echo "done"
			if [ $(numContainers) -ge 1 ] ; then
				echo "Some containers still running. Check and re-run command if needed."
			fi
		fi
	
		docker rmi registry:kube-single-node	
		docker rmi etcd:kube-single-node
		docker rmi hyperkube:kube-single-node
	fi
elif [ "$1x" = "startx" ] ; then
	if [ $(numContainers) -ge 1 ] ; then
		echo "There are containers running on the system. Clean-up and re-run." >&2
		exit 1
	fi

	echo -n "Starting containerized single node Kubernetes... "

	export REGISTRY_IMAGE=docker.io/registry:2
	export ETCD_IMAGE=gcr.io/google_containers/etcd:2.2.1
	export K8S_IMAGE=gcr.io/google_containers/hyperkube-amd64:v1.2.1
	export GOOS=linux
	export GOARCH=amd64

	docker pull ${REGISTRY_IMAGE}
	docker tag  ${REGISTRY_IMAGE} registry:kube-single-node

	docker pull ${ETCD_IMAGE}
	docker tag  ${ETCD_IMAGE} etcd:kube-single-node

	docker pull ${K8S_IMAGE}
	docker tag  ${K8S_IMAGE} hyperkube:kube-single-node

	manifestsVolume=${PWD}/$(dirname $0)"/etc/kubernetes/manifests"
	
	docker run -d \
	    --volume=/:/rootfs:ro \
	    --volume=/sys:/sys:rw \
	    --volume=/var/lib/docker/:/var/lib/docker:rw \
	    --volume=/var/lib/kubelet/:/var/lib/kubelet:rw \
	    --volume=/var/run:/var/run:rw \
	    --volume=${manifestsVolume}:/etc/kubernetes/manifests:ro \
	    --net=host \
	    --pid=host \
	    --privileged \
	    hyperkube:kube-single-node \
	    /hyperkube kubelet \
	        --containerized \
	        --hostname-override=127.0.0.1 \
	        --api-servers=http://localhost:8080 \
	        --config=/etc/kubernetes/manifests \
	        --allow-privileged --v=2 \
		>/dev/null 2>&1

	if [ $? -ne 0 ] ; then
		echo "failed!"
		docker rmi registry:kube-single-node	
		docker rmi etcd:kube-single-node
		docker rmi hyperkube:kube-single-node
		exit 1
	fi

	echo "done"
	
	if [ -d $HOME/bin ] ; then
		if [ ! -x $HOME/bin/kubectl ] ; then
			curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/${GOOS}/${GOARCH}/kubectl > $HOME/bin/kubectl
			chmod +x $HOME/bin/kubectl
		fi
	fi
else
	usage
	exit 1
fi

