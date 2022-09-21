#!/bin/bash

CLOUD_PROVIDER=${1:-$CLOUD_PROVIDER}
CRI_RUNTIME_ENDPOINT=${CRI_RUNTIME_ENDPOINT:-/run/peerpod/cri-runtime.sock}
optionals+=""

if [[ -S ${CRI_RUNTIME_ENDPOINT} ]]; then # will skip if socket isn't exist in the container
	optionals+="-cri-runtime-endpoint ${CRI_RUNTIME_ENDPOINT} "
fi


test_vars() {
        for i in $@; do
                [ -z ${!i} ] && echo "\$$i is NOT set" && EXT=1
        done
        [[ -n $EXT ]] && exit 1
}

aws() {
set -x

if [[ "${PODVM_LAUNCHTEMPLATE_NAME}" ]]; then
	optionals+="-use-lt -aws-lt-name ${PODVM_LAUNCHTEMPLATE_NAME}"
else
	optionals+="-imageid ${PODVM_AMI_ID} "
	optionals+="-instance-type ${PODVM_INSTANCE_TYPE:-t3.small} "
	optionals+="-securitygroupid ${AWS_SG_ID} "
	optionals+="-keyname ${SSH_KP_NAME} "
	optionals+="-subnetid ${AWS_SUBNET_ID} "
fi

cloud-api-adaptor-aws aws \
	-aws-access-key-id ${AWS_ACCESS_KEY_ID} \
	-aws-secret-key ${AWS_SECRET_ACCESS_KEY} \
	-aws-region ${AWS_REGION} \
	-pods-dir /run/peerpod/pods \
	${optionals} \
	-socket /run/peerpod/hypervisor.sock
}

libvirt() {
test_vars LIBVIRT_URI
set -x
cloud-api-adaptor-libvirt libvirt \
	-uri ${LIBVIRT_URI} \
	-data-dir /opt/data-dir \
	-pods-dir /run/peerpod/pods \
	-network-name ${LIBVIRT_NET:-default} \
	-pool-name ${LIBVIRT_POOL:-default} \
	${optionals} \
	-socket /run/peerpod/hypervisor.sock
}

vsphere() {
set -x
cloud-api-adaptor-vsphere vsphere \
        -vcenter-url ${GOVC_URL}  \
	-user-name ${GOVC_USERNAME} \
	-password ${GOVC_PASSWORD} \
	-data-store datastore2 \
	-pods-dir /run/peerpod/pods \
	-deploy-folder peerpods \
	-template podvm-new-template \
	${optionals} \
	-socket /run/peerpod/hypervisor.sock
}

help_msg() {
	cat <<EOF
Usage:
	CLOUD_PROVIDER=aws|libvirt|vsphere $0
or
	$0 aws|libvirt|vsphere
in addition all cloud provider specific env variables must be set and valid
(CLOUD_PROVIDER is currently set to "$CLOUD_PROVIDER")
EOF
}

if [[ "$CLOUD_PROVIDER" == "aws" ]]; then
	aws
elif [[ "$CLOUD_PROVIDER" == "libvirt" ]]; then
	libvirt
elif [[ "$CLOUD_PROVIDER" == "vsphere" ]]; then
	vsphere
else
	help_msg
fi