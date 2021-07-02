# Terraform Provider k8sbootstrap

After setting up a fresh Kubernetes cluster you need to get the CA certificate to connect to it.
This provider can be used for fetching the CA cert and preparing a kubeconfig.

Run the following command to build the provider

```shell
go build -o terraform-provider-k8sbootstrap
```
