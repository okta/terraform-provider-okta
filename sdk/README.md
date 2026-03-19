# LOCAL SDK

This is a local copy of the original v2 okta-sdk-golang and should be treated
as read-only. If the local SDK is missing a needed API call look for that
behavior in the [okta-sdk-golang v3](https://github.com/okta/okta-sdk-golang)
client.

Resources and data sources making use of the Terraform SDK v2 access the
[okta-sdk-golang v3](https://github.com/okta/okta-sdk-golang) as
`getOktaV3ClientFromMetadata(m)`. 

Resources and data sources making use of the Terraform Plugin Framework access
the [okta-sdk-golang v3](https://github.com/okta/okta-sdk-golang) as
`r.oktaSDKClientV3`.

NOTE: Any new resources or data sources contributed to this project should be
implemented with the [Terraform Plugin
Framework](https://developer.hashicorp.com/terraform/plugin/framework) and not
the [Terraform Plugin
SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2).
