![Build Status](https://github.com/okta/terraform-provider-okta/actions/workflows/release.yml/badge.svg)
<br/><br/>

<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-text.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>

<a href="https://www.okta.com/">
    <img src="https://www.okta.com/sites/default/files/Dev_Logo-03_Large.png" alt="OKTA logo" title="OKTA" height="50" />
</a>

# Terraform Provider for Okta

The Terraform Okta provider is a plugin for Terraform that allows for the full lifecycle management of Okta resources.
This provider is maintained internally by the Okta development team.

## Examples

All the resources and data sources has [one or more examples](./examples) to give you an idea of how to use this
provider to build your own Okta infrastructure. Provider's official documentation is located in the
[official terraform registry](https://registry.terraform.io/providers/okta/okta/latest/docs), or [here](./website/docs)
in form of raw markdown files.

# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.14.0 or newer (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please
check the [requirements](#requirements) before proceeding).

_Note:_ This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside
your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your
home directory outside the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:okta/terraform-provider-okta.git
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-okta
...
```

## Testing the Provider

In order to test the provider, you can run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run. Please
read [Running an Acceptance Test](https://github.com/okta/terraform-provider-okta/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test)
in the contribution guidelines for more information on usage.

```sh
$ make testacc
```

## Using the Provider

To use a released provider in your Terraform environment,
run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see
the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions)
.

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build
instructions above), follow the instructions
to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing the
custom-built provider into your plugins' directory, run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on
the [provider's website](https://registry.terraform.io/providers/okta/okta/latest/docs).

## Contributing

Terraform is the work of thousands of contributors. We really appreciate your help!

To contribute, please read the contribution
guidelines: [Contributing to Terraform - Okta Provider](.github/CONTRIBUTING.md)

Issues on GitHub are intended to be related to the bugs or feature requests with provider codebase.
See [Plugin SDK Community](https://www.terraform.io/community)
and [Discuss forum](https://discuss.hashicorp.com/c/terraform-providers/31/none) for a list of community resources to
ask questions about Terraform.
