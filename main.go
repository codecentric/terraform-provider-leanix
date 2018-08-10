package leanix

import (
	"terraform-provider-leanix/leanix"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return leanix.Provider()
		},
	})
}
