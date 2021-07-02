package bootstrap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func dataSourceKubeconfig() *schema.Resource {
	return &schema.Resource{
		Description: "Build kubeconfig from inputs - no verification done",
		ReadContext: dataSourceKubeconfigRead,
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Description: "Server URL to connect to (e.g. https://localhost:6443)",
				Optional:    true,
				Default:     "https://localhost:6443",
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "Namespace to use for operations",
				Optional:    true,
				Default:     "default",
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Bearer Token to use for authentication",
				Required:    true,
				Sensitive:   true,
			},
			"ca_crt": {
				Type:        schema.TypeString,
				Description: "Server CA certificate",
				Required:    true,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "Kubeconfig",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func getKubeconfig(server, caCrt, namespace, token string) (string, error) {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   server,
		CertificateAuthorityData: []byte(caCrt),
	}

	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:   "default-cluster",
		Namespace: namespace,
		AuthInfo:  namespace,
	}

	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos[namespace] = &clientcmdapi.AuthInfo{
		Token: token,
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	b, err := clientcmd.Write(clientConfig)
	return string(b), err
}

func dataSourceKubeconfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	kubeconfig, err := getKubeconfig(
		d.Get("server").(string),
		d.Get("ca_crt").(string),
		d.Get("namespace").(string),
		d.Get("token").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("kubeconfig", kubeconfig); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", d.Get("server").(string), d.Get("namespace").(string)))

	return diags
}
