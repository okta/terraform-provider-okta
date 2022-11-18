---
layout: "okta"
page_title: "Provider: Okta"
sidebar_current: "docs-okta-index"
description: |-
  The Okta provider is used to interact with the resources supported by Okta. The provider needs to be configured with the proper credentials before it can be used.
---

# Okta Provider

The Okta provider is used to interact with the resources supported by Okta. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources and data sources.

In case the provider configuration is still using old `"oktadeveloper/okta"` source, please change it to `"okta/okta"` and run
`terraform state replace-provider oktadeveloper/okta okta/okta`. Okta no longer supports `"oktadeveloper/okta"`.

## Example Usage

Terraform 0.14 and later:

```hcl
terraform {
  required_providers {
    okta = {
      source = "okta/okta"
      version = "~> 3.39"
    }
  }
}

# Configure the Okta Provider
provider "okta" {
  org_name  = "dev-123456"
  base_url  = "oktapreview.com"
  api_token = "xxxx"
}
```

For the resources and data sources examples, please check the [examples](https://github.com/okta/terraform-provider-okta/tree/master/examples) directory.

## Authentication

The Okta provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Environment variables
- Provider Config

### Environment variables

You can provide your credentials via the `OKTA_ORG_NAME`, `OKTA_BASE_URL`, `OKTA_ACCESS_TOKEN`, `OKTA_API_TOKEN`,
`OKTA_API_CLIENT_ID`, `OKTA_API_SCOPES` and `OKTA_API_PRIVATE_KEY` environment variables, representing your Okta
Organization Name, Okta Base URL (i.e. `"okta.com"` or `"oktapreview.com"`), Okta Access Token, Okta API Token,
Okta Client ID, Okta API scopes and Okta API private key respectively.

```hcl
provider "okta" {}
```

Usage:

```sh
$ export OKTA_ORG_NAME="dev-123456"
$ export OKTA_BASE_URL="oktapreview.com"
$ export OKTA_API_TOKEN="xxxx"
$ terraform plan
```

## Argument Reference

Note: `api_token` is mutually exclusive of the set `access_token`, `client_id`, `private_key`, and `scopes`. `api_token` is utilized for Okta's [SSWS Authorization Scheme](https://developer.okta.com/docs/reference/core-okta-api/#authentication) and applies to org level operations. `client_id`, `private_key`, and `scopes` are for [OAuth 2.0 client](https://developer.okta.com/docs/reference/api/apps/#add-oauth-2-0-client-application) authentication for application operations. `access_token` is used in situations where the caller has already performed the OAuth 2.0 client authentication process.

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Okta `provider` block:

- `org_name` - (Optional) This is the org name of your Okta account, for example `dev-123456.oktapreview.com` would have an org name of `dev-123456`. It must be provided, but it can also be sourced from the `OKTA_ORG_NAME` environment variable.

- `base_url` - (Optional) This is the domain of your Okta account, for example `dev-123456.oktapreview.com` would have a base url of `oktapreview.com`. It must be provided, but it can also be sourced from the `OKTA_BASE_URL` environment variable.

- `http_proxy` - (Optional) This is a custom URL endpoint that can be used for unit testing or local caching proxies. Can also be sourced from the `OKTA_HTTP_PROXY` environment variable.

- `access_token` - (Optional) This is an OAuth 2.0 access token to interact with your Okta org. It can be sourced from the `OKTA_ACCESS_TOKEN` environment variable. `access_token` conflicts with `api_token`, `client_id`, `scopes` and `private_key`.

- `api_token` - (Optional) This is the API token to interact with your Okta org. It can also be sourced from the `OKTA_API_TOKEN` environment variable. `api_token` conflicts with `access_token`, `client_id`, `scopes` and `private_key`.

- `client_id` - (Optional) This is the client ID for obtaining the API token. It can also be sourced from the `OKTA_API_CLIENT_ID` environment variable. `client_id` conflicts with `access_token` and `api_token`.

- `scopes` - (Optional) These are scopes for obtaining the API token in form of a comma separated list. It can also be sourced from the `OKTA_API_SCOPES` environment variable. `scopes` conflicts with `access_token` and `api_token`.

- `private_key` - (Optional) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `OKTA_API_PRIVATE_KEY` environment variable. `private_key` conflicts with `access_token` and `api_token`.

- `private_key_id` - (Optional) This is the private key ID (kid) for obtaining the API token. It can also be sourced from `OKTA_API_PRIVATE_KEY_ID` environmental variable. `private_key_id` conflicts with `api_token`.

- `backoff` - (Optional) Whether to use exponential back off strategy for rate limits, the default is `true`.

- `min_wait_seconds` - (Optional) Minimum seconds to wait when rate limit is hit, the default is `30`.

- `max_wait_seconds` - (Optional) Maximum seconds to wait when rate limit is hit, the default is `300`.

- `max_retries` - (Optional) Maximum number of retries to attempt before returning an error, the default is `5`.

- `request_timeout` - (Optional) Timeout for single request (in seconds) which is made to Okta, the default is `0` (means no limit is set). The maximum value can be `300`.

- `max_api_capacity` - (Optional, experimental) sets what percentage of capacity the provider can use of the total
  rate limit capacity while making calls to the Okta management API endpoints. Okta API operates in one minute buckets.
  See Okta Management API Rate Limits: https://developer.okta.com/docs/reference/rl-global-mgmt. Can be set to a value between 1 and 100.
