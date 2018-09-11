package leanix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestLeanixWebhookSubscription_basic(t *testing.T) {
	var subscription WebhookSubscription

	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform/helper/acctest
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExampleResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccExampleResource(resourceName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the subscription object
					testAccCheckExampleResourceExists("leanix_webhook_subscription.test", &subscription),
					// verify remote values
					testAccCheckExampleSubscriptionValues(&subscription, resourceName),
					// verify local values https://www.terraform.io/docs/extend/testing/acceptance-tests/teststep.html#builtin-check-functions
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "identifier", resourceName),
				),
			},
		},
	})
}

func testAccCheckExampleSubscriptionValues(subscription *WebhookSubscription, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if subscription.Identifier != identifier {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", identifier, subscription.Identifier)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching subscription.
func testAccCheckExampleResourceExists(target string, subscription *WebhookSubscription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[target]
		if !ok {
			return fmt.Errorf("Not found: %s", target)
		}

		// retrieve the configured client from the test setup
		leanix := testAccProvider.Meta().(*LeanixClient)
		resp, err := leanix.ReadWebhookSubscription(rs.Primary.ID)

		if err != nil {
			return err
		}

		// If no error, assign the response subscription attribute to the subscription pointer
		*subscription = *resp
		return nil
	}
}

// testAccCheckExampleResourceDestroy verifies the subscription has been destroyed
func testAccCheckExampleResourceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	leanix := testAccProvider.Meta().(*LeanixClient)

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "example_widget" {
			continue
		}

		_, err := leanix.ReadWebhookSubscription(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Widget (%s) still exists.", rs.Primary.ID)
		}

		// If the error is equivelent to 404 not found, the subscription is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "Widget not found") {
			return err
		}
	}
	return nil
}

// testAccExampleResource returns an configuration for an Example Widget with the provided name
func testAccExampleResource(name string) string {
	return fmt.Sprintf(`
resource "leanix_webhook_subscription" "test" {
  identifier           = "%s"
  ignore_error         = false
  target_url           = "http://localhost:1234"
  target_method        = "POST"
  authorization_header = "Basic ${base64encode("user:pass")}"
  callback             = "throw delivery.payload;"
  workspace_constraint = "ANY"
  payload_mode         = "WRAPPED_EVENT"
  active               = true
  workspace_id         = "8751abbf-8093-410d-a090-10c7735952cf"

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_UPDATED"]
  }

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_ARCHIVED"]
  }
}`, name)
}
func TestPackageTagSets(t *testing.T) {
	input := [][]string{
		[]string{"a", "b"},
		[]string{"c", "d"},
	}
	expectedOutput := []map[string][]string{
		map[string][]string{
			"tags": []string{"a", "b"},
		},
		map[string][]string{
			"tags": []string{"c", "d"},
		},
	}

	actualOutput := packageTagSets(input)
	assertEqual(t, actualOutput, expectedOutput)
}
