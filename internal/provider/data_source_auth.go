package provider

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/avast/retry-go"
)

func dataSourceAuth() *schema.Resource {
	return &schema.Resource{
		Description: "Fetch CA certificate and build kubeconfig",
		ReadContext: dataSourceAuthRead,
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Description: "Server URL to connect to (e.g. https://localhost:6443)",
				Optional:    true,
				Default:     "https://localhost:6443",
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "Namespace to retrieve the k8s API certificate secret from",
				Optional:    true,
				Default:     "kube-system",
			},
			"secret_name": {
				Type:        schema.TypeString,
				Description: "Name of the secret containing the k8s API certificate",
				Optional:    true,
				Default:     "k3s-serving",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Description: "Accept any server certificate - probably the purpose in using this data source",
				Optional:    true,
				Default:     true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Bearer Token to use for authentication",
				Required:    true,
				Sensitive:   true,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Description: "Timeout in seconds to wait for the API to be responding",
				Optional:    true,
				Default:     300,
			},
			"ca_crt": {
				Type:        schema.TypeString,
				Description: "Server CA certificate to use for further requests",
				Computed:    true,
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

func dataSourceAuthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: d.Get("server").(string)}},
	).ClientConfig()
	if err != nil {
		return diag.FromErr(err)
	}
	config.Insecure = d.Get("insecure").(bool)
	config.BearerToken = d.Get("token").(string)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return diag.FromErr(err)
	}

	cty, cancel := context.WithTimeout(ctx, time.Duration(d.Get("timeout").(int))*time.Second)
	defer cancel()

	err = retry.Do(
		func() error {
			secret, err := clientset.CoreV1().Secrets(d.Get("namespace").(string)).Get(ctx, d.Get("secret_name").(string), metav1.GetOptions{})
			if err != nil {
				return err
			}
			if tlsCrtBytes, ok := secret.Data["tls.crt"]; ok {
				for b, r := pem.Decode(tlsCrtBytes); b != nil; b, r = pem.Decode(r) {
					cert, err := x509.ParseCertificate(b.Bytes)
					if err != nil {
						return err
					}
					if cert.IsCA {
						if err := d.Set("ca_crt", string(pem.EncodeToMemory(b))); err != nil {
							return err
						}
						d.SetId(secret.Name)
						return nil
					}
				}
				return fmt.Errorf("secret does not contain a CA certificate at key tls.crt")
			}
			return fmt.Errorf("secret does not have a key named tls.crt")
		},
		retry.Context(cty),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(500*time.Millisecond),
		retry.Attempts(uint(d.Get("timeout").(int)*2)),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("[DEBUG] Connection attempt %d: %v - retrying...", n, err)
		}),
	)
	if err != nil {
		return diag.FromErr(err)
	}

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

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
