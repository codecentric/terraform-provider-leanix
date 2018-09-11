package leanix

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEANIX_URL", "https://svc.leanix.net"),
				Description: "LeanIX service URL.",
			},
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEANIX_API_TOKEN", nil),
				Description: "The LeanIX API token required to authenticate with LeanIX. See https://dev.leanix.net/docs/authentication for details.",
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
		d.Get("api_token").(string),
	), nil
}
