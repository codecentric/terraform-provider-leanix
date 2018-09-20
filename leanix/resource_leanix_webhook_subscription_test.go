package leanix

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// https://www.terraform.io/docs/extend/best-practices/testing.html
func TestLeanixWebhookSubscription_basic(t *testing.T) {
	var subscription WebhookSubscription

	// generate a random name for each subscription test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform/helper/acctest
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testSubscriptionResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testSubscriptionResource(resourceName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the subscription object
					testCheckResourceExists("leanix_webhook_subscription.test", &subscription),
					// verify remote values
					testCheckSubscriptionResourceValues(&subscription, resourceName),
					// verify local values https://www.terraform.io/docs/extend/testing/acceptance-tests/teststep.html#builtin-check-functions
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "identifier", resourceName),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "ignore_error", "false"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "target_url", "http://localhost:1234"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "target_method", "POST"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "authorization_header", "Basic "+base64.StdEncoding.EncodeToString([]byte("user:pass"))),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "callback", "throw delivery.payload;"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "workspace_constraint", "ANY"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "payload_mode", "WRAPPED_EVENT"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "active", "true"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "workspace_id", "8751abbf-8093-410d-a090-10c7735952cf"),
					resource.TestCheckResourceAttr("leanix_webhook_subscription.test", "tag_set.#", "2"),
				),
			},
		},
	})
}

func testCheckSubscriptionResourceValues(subscription *WebhookSubscription, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var expectedIdentifier = resourceName
		if subscription.Identifier != expectedIdentifier {
			return fmt.Errorf("Bad subscription.Identifier, expected \"%s\", got: %#v", expectedIdentifier, subscription.Identifier)
		}
		var expectedIgnoreError = false
		if subscription.IgnoreError != expectedIgnoreError {
			return fmt.Errorf("Bad subscription.IgnoreError, expected \"%v\", got: %#v", expectedIgnoreError, subscription.IgnoreError)
		}
		var expectedTargetUrl = "http://localhost:1234"
		if subscription.TargetUrl != expectedTargetUrl {
			return fmt.Errorf("Bad subscription.TargetUrl, expected \"%s\", got: %#v", expectedTargetUrl, subscription.TargetUrl)
		}
		var expectedTargetMethod = "POST"
		if subscription.TargetMethod != expectedTargetMethod {
			return fmt.Errorf("Bad subscription.TargetMethod, expected \"%s\", got: %#v", expectedTargetMethod, subscription.TargetMethod)
		}
		var expectedAuthorizationHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
		if subscription.AuthorizationHeader != expectedAuthorizationHeader {
			return fmt.Errorf("Bad subscription.AuthorizationHeader, expected \"%s\", got: %#v", expectedAuthorizationHeader, subscription.AuthorizationHeader)
		}
		var expectedCallback = "throw delivery.payload;"
		if subscription.Callback != expectedCallback {
			return fmt.Errorf("Bad subscription.Callback, expected \"%s\", got: %#v", expectedCallback, subscription.Callback)
		}
		var expectedWorkspaceConstraint = "ANY"
		if subscription.WorkspaceConstraint != expectedWorkspaceConstraint {
			return fmt.Errorf("Bad subscription.WorkspaceConstraint, expected \"%s\", got: %#v", expectedWorkspaceConstraint, subscription.WorkspaceConstraint)
		}
		var expectedPayloadMode = "WRAPPED_EVENT"
		if subscription.PayloadMode != expectedPayloadMode {
			return fmt.Errorf("Bad subscription.PayloadMode, expected \"%s\", got: %#v", expectedPayloadMode, subscription.PayloadMode)
		}
		var expectedActive = true
		if subscription.Active != expectedActive {
			return fmt.Errorf("Bad subscription.Active, expected \"%v\", got: %#v", expectedActive, subscription.Active)
		}
		var expectedWorkspaceId = "8751abbf-8093-410d-a090-10c7735952cf"
		if subscription.WorkspaceId != expectedWorkspaceId {
			return fmt.Errorf("Bad subscription.WorkspaceId, expected \"%s\", got: %#v", expectedWorkspaceId, subscription.WorkspaceId)
		}
		var expectedTagSets = [][]string{
			[]string{"pathfinder", "FACT_SHEET_UPDATED"},
			[]string{"pathfinder", "FACT_SHEET_ARCHIVED"},
		}
		if !reflect.DeepEqual(subscription.TagSets, expectedTagSets) {
			return fmt.Errorf("Bad subscription.TagSets, expected \"%s\", got: %#v", expectedTagSets, subscription.TagSets)
		}
		return nil
	}
}

// testCheckResourceExists queries the API and retrieves the matching subscription.
func testCheckResourceExists(target string, subscription *WebhookSubscription) resource.TestCheckFunc {
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

// testSubscriptionResourceDestroy verifies the subscription has been destroyed
func testSubscriptionResourceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	leanix := testAccProvider.Meta().(*LeanixClient)

	// loop through the resources in state, verifying each subscription is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "leanix_webhook_subscription" {
			continue
		}

		_, err := leanix.ReadWebhookSubscription(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Subscription (%s) still exists.", rs.Primary.ID)
		}

		// If the error is equivelent to 404 not found, the subscription is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "No such subscription") {
			return err
		}
	}
	return nil
}

// testSubscriptionResource returns an configuration for an example subscription with the provided name
func testSubscriptionResource(name string) string {
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
