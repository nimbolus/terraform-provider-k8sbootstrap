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
