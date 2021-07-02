package bootstrap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"k8sbootstrap_auth":       dataSourceAuth(),
			"k8sbootstrap_kubeconfig": dataSourceKubeconfig(),
		},
	}
}
