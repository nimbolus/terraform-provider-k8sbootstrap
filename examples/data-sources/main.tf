terraform {
  required_providers {
    k8sbootstrap = {
      source = "github.com/nimbolus/k8sbootstrap"
    }
  }
}

variable "server" {
  type    = string
  default = "https://k8s.example.com:6443"
}

variable "token" {
  type = string
}

data "k8sbootstrap_auth" "ca" {
  server = var.server
  token  = var.token
}

output "ca_crt" {
  value = data.k8sbootstrap_auth.ca.ca_crt
}
output "kubeconfig_real" {
  value     = data.k8sbootstrap_auth.ca.kubeconfig
  sensitive = true
}

data "k8sbootstrap_kubeconfig" "templated" {
  server = "https://not.reachable.example.com"
  token  = "not-a-real-token"
  ca_crt = "no real ca crt"
}
output "kubeconfig_templated" {
  value     = data.k8sbootstrap_kubeconfig.templated.kubeconfig
  sensitive = true
}
