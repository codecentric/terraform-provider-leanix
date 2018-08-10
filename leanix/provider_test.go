package leanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = new(schema.Provider)
}
