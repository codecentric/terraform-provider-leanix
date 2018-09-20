package leanix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEANIX_URL", "https://svc.leanix.net"),
				Description: "LeanIX service URL.",
			},
			"auth_header": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEANIX_AUTH_HEADER", nil),
				Description: "The LeanIX authentication header based on API token or client secret required to authenticate with LeanIX. See https://dev.leanix.net/docs/authentication for details.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"leanix_webhook_subscription": resourceLeanixWebhookSubscription(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return NewLeanixClient(
		d.Get("url").(string),
		d.Get("auth_header").(string),
	), nil
}
