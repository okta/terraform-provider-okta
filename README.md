# terraform-provider-okta
terraform provider for okta authentication service

## Prerequisites
1. a golang environment setup (golang.org/doc/install) + terraform installed
2. an okta organization account (development account signup: developer.okta.com/signup)
3. an api token (developer.okta.com/docs/api/getting_started/getting_a_token.html)
4. you *should* create a service user with superuser privileges and create the api token while logged in as that user

## Building the provider
```
go get
```
```
go build -o .terraform/plugins/linux_amd64/terraform-provider-okta
```
```
terraform init -plugin-dir=.terraform/plugins/linux_amd64
```
