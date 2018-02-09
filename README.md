Terraform Provider Okta
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Terraform team at [Articulate](https://articulate.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Usage
---------------------

This plugin requires two inputs to run: the okta organization name and the okta api token. The okta base url is not required and will default to "okta.com" if left out.

You can specify the inputs in your tf plan:

```
provider "okta" {
  org_name  = <okta instance name, e.g. dev-151148>
  api_token = <okta instance api token with the Administrator role>
  base_url  = <okta base url, e.g. oktapreview.com>
}
```

OR you can specify environment variables:

```
OKTA_ORG_NAME=<okta instance name, e.g. dev-151148>
OKTA_API_TOKEN=<okta instance api token with the Administrator role>
OKTA_BASE_URL=<okta base url, e.g. oktapreview.com>
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-okta`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-okta
```

Enter the provider directory and build the provider
This provider has a dependency on the [Go Okta SDK](http://github.com/articulate/oktasdk-go)

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-okta
$ go get -v
$ make build
```

For local development, I've found the below commands helpful. Run them from inside the terraform-provider-okta directory

```sh
$ go build -o .terraform/plugins/$GOOS_$GOARCH/terraform-provider-okta
$ terraform init -plugin-dir=.terraform/plugins/$GOOS_$GOARCH
```

Using the provider
----------------------

Example terraform plan:

```
provider "okta" {
  org_name  = "dev-XXXXX"
  api_token = "XXXXXXXXXXXXXXXXXXXXXXXX"
  base_url  = "oktapreview.com"
}

resource "okta_users" "blah" {
  firstname = "blah"
  lastname  = "blergh"
  email     = "XXXXX@XXXXXXXX.XXX"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-okta
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
