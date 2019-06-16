# Development

This doc explains the development workflow so you can get started [contributing](CONTRIBUTING.md) to the Okta Terraform Provider!

## Getting started

First you will need to setup your GitHub account and create a fork:

1. Create [a GitHub account](https://github.com/join)
1. Setup [GitHub access via
   SSH](https://help.github.com/articles/connecting-to-github-with-ssh/)
1. [Create and checkout a repo fork](#checkout-your-fork)

Once you have those, you can iterate on the provider:

1. [Run tests](#testing)

When you're ready, you can [create a PR](#creating-a-pr)!

## Checkout your fork

To check out this repository:

1. Create your own [fork of this
  repo](https://help.github.com/articles/fork-a-repo/)
2. Clone it to your machine:

  ```shell
  mkdir -p ${GOPATH}/src/github.com/articulate
  cd ${GOPATH}/src/github.com/articulate
  git clone git@github.com:${YOUR_GITHUB_USERNAME}/terraform-provider-okta.git
  cd terraform-provider-okta
  git remote add upstream git@github.com:articulate/terraform-provider-okta.git
  git remote set-url --push upstream no_push
  ```

_Adding the `upstream` remote sets you up nicely for regularly [syncing your
fork](https://help.github.com/articles/syncing-a-fork/)._

## Testing

The provider has both [unit tests](#unit-tests) and [acceptance tests](#acceptance-tests).

### Unit Tests

To run unit tests simply run

```shell
make test
```

### Acceptance tests

Acceptance tests are run against real infrastructure and thus require credentials for an Okta org. Be sure not to run this anywhere but a test org. Start by copying `.env.sample` and installing a version of `dotenv` cli or simply export these values into your environment.

```shell
cp .env.sample .env
dotenv make testacc
```

## Creating a PR

When you have changes you would like to propose to kritis, you will need to:

1. Ensure the commit message(s) describe what issue you are fixing and how you are fixing it
   (include references to [issue numbers](https://help.github.com/articles/closing-issues-using-keywords/)
   if appropriate)
1. [Create a pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)
1. Post a screenshot of the passing ACC tests. If you do not have access to an Okta org, you can request a maintainer run the ACC tests. These tests make ALOT of API calls so free dev accounts, the configured number of retries, and the backoff duration may not be enough to get through all of the tests.

### Reviews

Each PR must be reviewed by a maintainer. In order for a PR to merge you must post a screenshot of the ACC tests passing. We do not run these via Travis due to Okta rate limiting. It gets to be untenable. Maintainers have access to Okta orgs they can run this against so if you are an outside contributor feel free to request ACC tests be run.
