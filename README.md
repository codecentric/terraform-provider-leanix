# Terraform Provider LeanIX

[![Build Status](https://travis-ci.org/codecentric/terraform-provider-leanix.svg?branch=master)](https://travis-ci.org/codecentric/terraform-provider-leanix)

## Description

This is a custom [Terraform provider](https://www.terraform.io/docs/providers/index.html) to manage LeanIX resources.

## Provider Configuration

The LeanIX provider requires a valid LeanIX URL and an API key to authenticate. Please make sure that the API key has the required permissions to manage the resources you want to use. You can either set the URL and API token directly in the provider, or use the environment variables.

```hcl
provider "leanix" {
  url         = "https://svc.leanix.net"                                                     // = LEANIX_URL
  auth_header = "Basic ${base64encode("apitoken:aVQEzWKwE2sSp3rhVKWwaVQEzWKwE2sSp3rhVKWw")}" // = LEANIX_AUTH_HEADER
}
```

## Supported Resources

### Webhook Subscription

The webhook subscription resource is used to manage webhook subscriptions in LeanIX. For more information please refer to the official [webhooks documentation](https://dev.leanix.net/docs/webhooks) or [API reference](https://svc.leanix.net/services/webhooks/v1/docs/#/).

#### Example

```hcl
resource "leanix_webhook_subscription" "example" {
  identifier           = "mySubscription"
  ignore_error         = true
  target_url           = "https://company.domain/my-endpoint"
  target_method        = "POST"
  authorization_header = "Basic dXNlcjpwYXNzCg=="
  callback             = "delivery.payload = { payload: \"foo\" };"
  workspace_constraint = "ANY"
  payload_mode         = "WRAPPED_EVENT"
  active               = true
  workspace_id         = "aa32abbf-8093-410d-a090-10c7735952cf"

  tag_set {
    tag {
      value = "pathfinder"
    }
    tag {
      value = "FACT_SHEET_CREATED"
    }
  }

  tag_set {
    tag {
      value = "pathfinder"
    }
    tag {
      value = "FACT_SHEET_UPDATED"
    }
  }
}
```

## Installation

The provider needs to be installed as a [third-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins). A quick and convenient way is to download the correct binary for your platform from the [GitHub release page](https://github.com/codecentric/terraform-provider-leanix/releases).

Then you can rename it to `terraform-provider-leanix` and place it inside the folder containing your `.tf` files. Terraform will load the provider plugin on startup automatically based on the filename `terraform-provider-*`.

## Building from Source

1. Install dependencies with `go get`
2. Execute tests with `go test ./...`
3. Execute acceptance tests (optional)
   ```sh
   LEANIX_AUTH_HEADER="<leanix_auth_header>" \
   LEANIX_URL="<leanix_url>" \
   TF_ACC=1 \
   go test -v ./...
   ```
4. Package provider executable with `go build`

## Release

To release a new version of the provider you need to add a git tag in the form of `v${x}.${y}.\${z}`, e.g. `v1.2.3`. Pre-release tags are also supported (`v1.1.2-rc1`, `v2.0.0-alpha1`).

After pushing the tags the binaries will be built and published as a GitHub release by GitHub Actions.

```sh
git checkout master
git tag v2.8.1
git push origin master --tags
```
