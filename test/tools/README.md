# Introduction

tools for development and testing.

## provisioner-cli

`provisioner-cli` provides a cli program that leverage the [cluster provisioner package](../provisioner) to create VPC, Subnet and Cluster and other necessary resources, and then deploy the cloud-api-adaptor operator to enable the function in the created cluster. Which is also used to upload a VM image to cloud vendor.

### Build provisioner-cli
In the root directory of `test/tools`, run command as below to build the cli program:
```bash
make provisioner-cli
```

Program is generated: `test/tools/caa-provisioner-cli`

### Use provisioner-cli
In directory `test/tools`, run commands like:
```bash
export TEST_E2E_PODVM_IMAGE=${POD_IMAGE_FILE_PATH}
export LOG_LEVEL=${LOG_LEVEL}
export CLOUD_PROVIDER=${CLOUD_PROVIDER}
export TEST_E2E_PROVISION_FILE=${PROPERTIES_FILE_PATH}
export TEST_E2E_PROVISION="yes"
./caa-provisioner-cli -actions=${ACTION}
```
`ACTION` supports `provision`, `deprovision` and `uploadimage`.

### Add a new provider support
`ibmcloud`, `azure` and `libvirt` providers are supported now, to add a new provider please add it in [cluster provisioner package](../provisioner)