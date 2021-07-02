variable "token" {
  type = string
}

data "k8sbootstrap_kubeconfig" "templated" {
  server = "https://not.reachable.example.com"
  token  = var.token
  ca_crt = "put your ca cert here"
}

output "kubeconfig" {
  value     = data.k8sbootstrap_kubeconfig.templated.kubeconfig
  sensitive = true
}
