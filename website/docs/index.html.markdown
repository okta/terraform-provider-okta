---
layout: 'okta'
page_title: 'Provider: Okta'
sidebar_current: 'docs-okta-index'
description: |-
  The Okta provider is used to interact with the resources supported by Okta. The provider needs to be configured with the proper credentials before it can be used.
---

# Okta Provider

The Okta provider is used to interact with the resources supported by Okta. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    okta = {
      source = "okta/okta"
      version = "~> 3.6"
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

Terraform 0.12 and earlier:

```hcl

# Configure the Okta Provider
provider "okta" {
  org_name  = "dev-123456"
  base_url  = "oktapreview.com"
  api_token = "xxxx"
}
```

## Authentication

The Okta provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Environment variables
- Provider Config

### Environment variables

You can provide your credentials via the `OKTA_ORG_NAME`, `OKTA_BASE_URL`, `OKTA_API_TOKEN`, `OKTA_CLIENT_ID`, 
`OKTA_API_SCOPES` and `OKTA_PRIVATE_KEY` environment variables, representing your Okta Organization Name, 
Okta Base URL (i.e. `"okta.com"` or `"oktapreview.com"`), Okta API Token, Okta Client ID, Okta API scopes 
and Okta API private key respectively.

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

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Okta `provider` block:

- `org_name` - (Optional) This is the org name of your Okta account, for example `dev-123456.oktapreview.com` would have an org name of `dev-123456`. It must be provided, but it can also be sourced from the `OKTA_ORG_NAME` environment variable.

- `base_url` - (Optional) This is the domain of your Okta account, for example `dev-123456.oktapreview.com` would have a base url of `oktapreview.com`. It must be provided, but it can also be sourced from the `OKTA_BASE_URL` environment variable.

- `api_token` - (Optional) This is the API token to interact with your Okta org (either `"api_token"` or `"client_id"`, `"scopesv` and `"private_key"` must be provided). It can also be sourced from the `OKTA_API_TOKEN` environment variable.

- `client_id` - (Optional) This is the client ID for obtaining the API token. It can also be sourced from the `OKTA_CLIENT_ID` environment variable. 

- `scopes` - (Optional) These are scopes for obtaining the API token in form of a comma separated list. It can also be sourced from the `OKTA_API_SCOPES` environment variable.

- `private_key` - (Optional) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `OKTA_PRIVATE_KEY` environment variable.

- `backoff` - (Optional) Whether to use exponential back off strategy for rate limits, the default is `true`.

- `min_wait_seconds` - (Optional) Minimum seconds to wait when rate limit is hit, the default is `30`.

- `max_wait_seconds` - (Optional) Maximum seconds to wait when rate limit is hit, the default is `300`.

- `max_retries` - (Optional) Maximum number of retries to attempt before returning an error, the default is `5`.

- `request_timeout` - (Optional) Timeout for single request (in seconds) which is made to Okta, the default is `0` (means no limit is set). The maximum value can be `100`.
