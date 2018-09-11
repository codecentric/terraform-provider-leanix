package leanix

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = new(schema.Provider)
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"leanix": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("LEANIX_URL") == "" {
		t.Fatal("LEANIX_URL must be set for acceptance tests")
	}
	if os.Getenv("LEANIX_API_TOKEN") == "" {
		t.Fatal("LEANIX_API_TOKEN must be set for acceptance tests")
	}
}
