# Terraform Provider LeanIX

[![Build Status](https://travis-ci.org/codecentric/terraform-provider-leanix.svg?branch=master)](https://travis-ci.org/codecentric/terraform-provider-leanix)

## Description

This is a custom [Terraform provider](https://www.terraform.io/docs/providers/index.html) to manage LeanIX resources.

## Provider Configuration

The LeanIX provider requires a valid LeanIX URL and an API key to authenticate. Please make sure that the API key has the required permissions to manage the resources you want to use.

```
provider "leanix" {
  url       = "https://svc.leanix.net"
  api_token = "aVQEzWKwE2sSp3rhVKWwaVQEzWKwE2sSp3rhVKWw"
}
```

## Supported Resources

### Webhook Subscription

The webhook subscription resource is used to manage webhook subscriptions in LeanIX. For more information please refer to the official [webhooks documentation](https://dev.leanix.net/docs/webhooks) or [API reference](https://svc.leanix.net/services/webhooks/v1/docs/#/).

#### Example

```
resource "leanix_webhook_subscription" "example" {
  identifier           = "mySubscription"
  ignore_error         = true
  target_url           = "https://company.domain/my-endpoint"
  target_method        = "POST"
  authorization_header = "Basic dXNlcjpwYXNzCg=="
  callback             = "delivery.payload = { payload: \"foo\" };"
  workspace_constraint = "ANY"
  active               = true
  workspace_id         = "aa32abbf-8093-410d-a090-10c7735952cf"

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_CREATED"]
  }

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_UPDATED"]
  }

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_ARCHIVED"]
  }

  tag_set {
    tags = ["pathfinder", "FACT_SHEET_DELETED"]
  }
}
```

# Development

## Building from Source

1. Install dependencies with `go get`
2. Execute tests with `go test ./...`
2. Package provider executable with `go build`

# Installing the Provider

The provider needs to be installed as a [third-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins). A quick and convenient way is to copy the `terraform-provider-leanix` binary inside the folder containing your `.tf` files. Terraform will load the provider plugin on startup automatically based on the filename `terraform-provider-*`.