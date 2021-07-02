variable "server" {
  type    = string
  default = "https://k8s.example.com:6443"
}

variable "token" {
  type = string
}

data "k8sbootstrap_auth" "auth" {
  server = var.server
  token  = var.token
}

output "ca_crt" {
  value = data.k8sbootstrap_auth.auth.ca_crt
}
output "kubeconfig" {
  value     = data.k8sbootstrap_auth.auth.kubeconfig
  sensitive = true
}
