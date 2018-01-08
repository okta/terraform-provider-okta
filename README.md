# terraform-provider-okta
terraform provider for okta authentication service

## Prerequisites
1. an okta organization account (development account signup: developer.okta.com/signup)
2. an api token (developer.okta.com/docs/api/getting_started/getting_a_token.html)
3. you *should* create a service user with superuser privileges and create the api token logged in as that user

## Building the provider
```
go build -o .terraform/plugins/linux_amd64/terraform-provider-okta
```
terraform init -plugin-dir=.terraform/plugins/linux_amd64
```
