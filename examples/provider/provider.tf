terraform {
  required_providers {
    k8sbootstrap = {
      source = "nimbolus/k8sbootstrap"
    }
  }
}

provider "k8sbootstrap" {}
