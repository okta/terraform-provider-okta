# Contributing to Terraform - Okta Provider

**First:** if you're unsure or afraid of _anything_, ask for help! You can
submit a work in progress (WIP) pull request, or file an issue with the parts
you know. We'll do our best to guide you in the right direction, and let you
know if there are guidelines we will need to follow. We want people to be able
to participate without fear of doing the wrong thing.

Below are our expectations for contributors. Following these guidelines gives us
the best opportunity to work with you, by making sure we have the things we need
in order to make it happen. Doing your best to follow it will speed up our
ability to merge PRs and respond to the issues.

<!-- TOC depthFrom:2 -->

- [Issues](#issues)
  - [Issue Reporting Checklists](#issue-reporting-checklists)
    - [Bug Reports](#bug-reports)
    - [Feature Requests](#feature-requests)
    - [Questions](#questions)
  - [Issue Lifecycle](#issue-lifecycle)
- [Pull Requests](#pull-requests)
  - [Pull Request Lifecycle](#pull-request-lifecycle)
  - [Checklists for Contribution](#checklists-for-contribution)
    - [Documentation Update](#documentation-update)
    - [Enhancement/Bugfix to a Resource](#enhancementbugfix-to-a-resource)
    - [Adding Resource Import Support](#adding-resource-import-support)
    - [New Resource](#new-resource)
  - [Common Review Items](#common-review-items)
    - [Go Coding Style](#go-coding-style)
    - [Resource Contribution Guidelines](#resource-contribution-guidelines)
    - [Acceptance Testing Guidelines](#acceptance-testing-guidelines)
  - [Writing Acceptance Tests](#writing-acceptance-tests)
    - [Acceptance Tests Often Cost Money to Run](#acceptance-tests-often-cost-money-to-run)
    - [Acceptance Tests With VCR](#acceptance-tests-with-vcr)
    - [Running an Acceptance Test](#running-an-acceptance-test)
    - [Writing an Acceptance Test](#writing-an-acceptance-test)

<!-- /TOC -->

## Issues

### Issue Reporting Checklists

We welcome issues of all kinds including feature requests, bug reports, and
general questions. Below you'll find checklists with guidelines for well-formed
issues of each type.

#### [Bug Reports](https://github.com/okta/terraform-provider-okta/issues/new?template=Bug_Report.md)

- [ ] **Test against the latest release**: Make sure you test against the latest
      released version. It is possible we already fixed the bug you're experiencing.

- [ ] **Search for possible duplicate reports**: It's helpful to keep bug
      reports consolidated to one thread, so do a quick search on existing bug
      reports to check if anybody else has reported the same thing. You can [scope
      searches by the label "bug"](https://github.com/okta/terraform-provider-okta/issues?q=is%3Aopen+is%3Aissue+label%3Abug) to help narrow things down.

- [ ] **Include steps to reproduce**: Provide steps to reproduce the issue,
      along with your `.tf` files, with secrets removed, so we can try to
      reproduce it. Without this, it makes it much harder to fix the issue.

- [ ] **For panics, include `crash.log`**: If you experienced a panic, please
      create a [gist](https://gist.github.com) of the _entire_ generated crash log
      for us to look at. Double check no sensitive items were in the log.

#### [Feature Requests](https://github.com/okta/terraform-provider-okta/issues/new?labels=enhancement&template=Feature_Request.md)

- [ ] **Search for possible duplicate requests**: It's helpful to keep requests
      consolidated to one thread, so do a quick search on existing requests to
      check if anybody else has reported the same thing. You can [scope searches by
      the label "enhancement"](https://github.com/okta/terraform-provider-okta/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement) to help narrow things down.

- [ ] **Include a use case description**: In addition to describing the
      behavior of the feature you'd like to see added, it's helpful to also lay
      out the reason why the feature would be important and how it would benefit
      Terraform users.

#### [Questions](https://github.com/okta/terraform-provider-okta/issues/new?labels=question&template=Question.md)

- [ ] **Search for answers in Terraform documentation**: We're happy to answer
      questions in GitHub Issues, but it helps reduce issue churn and maintainer
      workload if you work to [find answers to common questions in the
      documentation](https://registry.terraform.io/providers/okta/okta/index.html). Oftentimes Question issues result in documentation updates
      to help future users, so if you don't find an answer, you can give us
      pointers for where you'd expect to see it in the docs.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a Terraform collaborator.
   Categorization is done via GitHub labels. We generally use a two-label
   system of (1) issue/PR type, and (2) section of the codebase. Type is
   one of "bug", "enhancement", "documentation", or "question", and section
   is usually the Okta service name.

3. An initial triage process determines whether the issue is critical and must
   be addressed immediately, or can be left open for community discussion.

4. The issue is addressed in a pull request or commit. The issue number will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.

5. The issue is closed. Sometimes, valid issues will be closed because they are
   tracked elsewhere or non-actionable. The issue is still indexed and
   available for future viewers, or can be re-opened if necessary.

## Pull Requests

We appreciate direct contributions to the provider codebase. Here's what to
expect:

- For pull requests that follow the guidelines, we will proceed to reviewing
  and merging, following the provider team's review schedule. There may be some
  internal or community discussion needed before we can complete this.
- Pull requests that don't follow the guidelines will be commented with what
  they're missing. The person who submits the pull request or another community
  member will need to address those requests before they move forward.

### Pull Request Lifecycle

1. [Fork the GitHub repository](https://help.github.com/en/articles/fork-a-repo), modify the code, 
   and [create a pull request](https://help.github.com/en/articles/creating-a-pull-request-from-a-fork).
   You are welcome to submit your pull request for the commentaries or review before
   it is fully completed by creating a [draft pull request](https://help.github.com/en/articles/about-pull-requests#draft-pull-requests)
   or adding `[WIP]` to the beginning of the pull request title.
   Please include specific questions or items you'd like feedback on.

1. Once you believe your pull request is ready to be reviewed, ensure the
   pull request is not a draft pull request by [marking it ready for review](https://help.github.com/en/articles/changing-the-stage-of-a-pull-request)
   or removing `[WIP]` from the pull request title if necessary, and a
   maintainer will review it. Follow [the checklists below](#checklists-for-contribution)
   to help ensure that your contribution can be easily reviewed and potentially
   merged.

1. One of Terraform's provider team members will look over your contribution and
   either approve it or provide comments letting you know if there is anything
   left to do. We do our best to keep up with the volume of PRs waiting for
   review, but it may take some time depending on the complexity of the work.

1. Once all outstanding comments and checklist items have been addressed, your
   contribution will be merged! Merged PRs will be included in the next
   Terraform release. The provider team takes care of updating the CHANGELOG as
   they merge.

1. In some cases, we might decide that a PR should be closed without merging.
   We'll make sure to provide clear reasoning when this happens.

### Checklists for Contribution

There are several kinds of contribution, each of which has its own
standards for a speedy review. The following sections describe guidelines for
each type of contribution.

#### Documentation Update

The [Terraform Okta Provider's website source][website] is in this repository
along with the code and tests. Below are some common items that will get
flagged during documentation reviews:

- [ ] **Reasoning for Change**: Documentation updates should include an explanation for why the update is needed.
- [ ] **Prefer Okta Documentation**: Documentation about Okta service features and valid argument values that are likely to update over time should link to Okta service user guides and API references where possible.
- [ ] **Large Example Configurations**: Example Terraform configuration that includes multiple resource definitions should be added to the repository `examples` directory instead of an individual resource documentation page. Each directory under `examples` should be self-contained to call `terraform apply` without special configuration.
- [ ] **Terraform Configuration Language Features**: Individual resource documentation pages and examples should refrain from highlighting particular Terraform configuration language syntax workarounds or features such as `variable`, `local`, `count`, and built-in functions.

#### Enhancement/Bugfix to a Resource

Working on existing resources is a great way to get started as a Terraform
contributor because you can work within existing code and tests to get a feel
for what to do.

In addition to the below checklist, please see the [Common Review
Items](#common-review-items) sections for more specific coding and testing
guidelines.

- [ ] **Acceptance test coverage of new behavior**: Existing resources each
      have a set of [acceptance tests][acctests] covering their functionality.
      These tests should exercise all the behavior of the resource. Whether you are
      adding something or fixing a bug, the idea is to have an acceptance test that
      fails if your code were to be removed. Sometimes it is sufficient to
      "enhance" an existing test by adding an assertion or tweaking the config
      that is used, but it's often better to add a new test. You can copy/paste an
      existing test and follow the conventions you see there, modifying the test
      to exercise the behavior of your code.
- [ ] **Documentation updates**: If your code makes any changes that need to
      be documented, you should include those doc updates in the same PR. This
      includes things like new resource attributes or changes in default values.
      The [Terraform website][website] source is in this repo and includes
      instructions for getting a local copy of the site up and running if you'd
      like to preview your changes.
- [ ] **Well-formed Code**: Do your best to follow existing conventions you
      see in the codebase, and ensure your code is formatted with `go fmt`. (The
      Travis CI build will fail if `go fmt` has not been run on incoming code.)
      The PR reviewers can help out on this front, and may provide comments with
      suggestions on how to improve the code.
- [ ] **Vendor additions**: Create a separate PR if you are updating the vendor
      folder. This is to avoid conflicts as the vendor versions tend to be fast-
      moving targets. We will plan to merge the PR with this change first.

#### Adding Resource Import Support

Adding import support for Terraform resources will allow existing infrastructure to be managed within Terraform. This type of enhancement generally requires a small to moderate amount of code changes.

Comprehensive code examples and information about resource import support can be found in the [Extending Terraform documentation](https://www.terraform.io/docs/extend/resources/import.html).

In addition to the below checklist and the items noted in the Extending Terraform documentation, please see the [Common Review Items](#common-review-items) sections for more specific coding and testing guidelines.

- [ ] _Resource Code Implementation_: In the resource code (e.g. `okta/resource_okta_service_thing.go`), implementation of `Importer` `State` function
- [ ] _Resource Acceptance Testing Implementation_: In the resource acceptance testing (e.g. `okta/resource_okta_service_thing_test.go`), implementation of `TestStep`s with `ImportState: true`
- [ ] _Resource Documentation Implementation_: In the resource documentation (e.g. `website/docs/r/service_thing.html.markdown`), addition of `Import` documentation section at the bottom of the page

#### New Resource

Implementing a new resource is a good way to learn more about how Terraform
interacts with upstream APIs. There are plenty of examples to draw from in the
existing resources, but you still get to implement something completely new.

In addition to the below checklist, please see the [Common Review
Items](#common-review-items) sections for more specific coding and testing
guidelines.

- [ ] **Minimal LOC**: It's difficult for both the reviewer and author to go
      through long feedback cycles on a big PR with many resources. We ask you to
      only submit **1 resource at a time**.
- [ ] **Acceptance tests**: New resources should include acceptance tests
      covering their behavior. See [Writing Acceptance
      Tests](#writing-acceptance-tests) below for a detailed guide on how to
      approach these.
- [ ] **Resource Naming**: Resources should be named `okta_<service>_<name>`,
      using underscores (`_`) as the separator. Resources are namespaced with the
      service name to allow easier searching of related resources, to align
      <!--the resource naming with the service for [Customizing Endpoints](https://www.terraform.io/docs/providers/okta/guides/custom-service-endpoints.html#available-endpoint-customizations),-->
      the resource naming with the service for Customizing Endpoints,
      and to prevent future conflicts with new Okta services/resources.
      For reference:

  - `service` is the Okta short service name that matches the entry in
    `endpointServiceNames` (created via the [New Service](#new-service)
    section)
  - `name` represents the conceptual infrastructure represented by the
    create, read, update, and delete methods of the service API. It should
    be a singular noun. For example, in an API that has methods such as
    `CreateThing`, `DeleteThing`, `DescribeThing`, and `ModifyThing` the name
    of the resource would end in `_thing`.

- [ ] **Arguments_and_Attributes**: The HCL for arguments and attributes should
      mimic the types and structs presented by the Okta API. API's arguments should be
      converted from `CamelCase` to `camel_case`.
- [ ] **Documentation**: Each resource gets a page in the Terraform
      documentation. The [Terraform website][website] source is in this
      repo and includes instructions for getting a local copy of the site up and
      running if you'd like to preview your changes. For a resource, you'll want
      to add a new file in the appropriate place and add a link to the sidebar for
      that page.
- [ ] **Well-formed Code**: Do your best to follow existing conventions you
      see in the codebase, and ensure your code is formatted with `go fmt`. (The
      Travis CI build will fail if `go fmt` has not been run on incoming code.)
      The PR reviewers can help out on this front, and may provide comments with
      suggestions on how to improve the code.
- [ ] **Vendor updates**: Create a separate PR if you are adding to the vendor
      folder. This is to avoid conflicts as the vendor versions tend to be fast-
      moving targets. We will plan to merge the PR with this change first.

### Common Review Items

The Terraform Okta Provider follows common practices to ensure consistent and
reliable implementations across all resources in the project. While there may be
older resource and testing code that predates these guidelines, new submissions
are generally expected to adhere to these items to maintain Terraform Provider
quality. For any guidelines listed, contributors are encouraged to ask any
questions and community reviewers are encouraged to provide review suggestions
based on these guidelines to speed up the review and merge process.

#### Go Coding Style

The following Go language resources provide common coding preferences that may be referenced during review, if not automatically handled by the project's linting tools.

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### Resource Contribution Guidelines

The following resource checks need to be addressed before your contribution can be merged. The exclusion of any applicable check may result in a delayed time to merge.

- [ ] **Passes Testing**: All code and documentation changes must pass unit testing, code linting, and website link testing. Resource code changes must pass all acceptance testing for the resource.
- [ ] **Avoids Optional and Required for Non-Configurable Attributes**: Resource schema definitions for read-only attributes should not include `Optional: true` or `Required: true`.
- [ ] **Avoids Resource Read Function in Data Source Read Function**: Data sources should fully implement their own resource `Read` functionality including duplicating `d.Set()` calls.
- [ ] **Avoids Reading Schema Structure in Resource Code**: The resource `Schema` should not be read in resource `Create`/`Read`/`Update`/`Delete` functions to perform looping or otherwise complex attribute logic. Use [`d.Get()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.Get) and [`d.Set()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.Set) directly with individual attributes instead.
- [ ] **Avoids ResourceData.GetOkExists()**: Resource logic should avoid using [`ResourceData.GetOkExists()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.GetOkExists) as its expected functionality is not guaranteed in all scenarios.
- [ ] **Implements Read After Create and Update**: Except where API eventual consistency prohibits immediate reading of resources or updated attributes, resource `Create` and `Update` functions should return the resource `Read` function.
- [ ] **Implements Immediate Resource ID Set During Create**: Immediately after calling the API creation function, the resource ID should be set with [`d.SetId()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.SetId) before other API operations or returning the `Read` function.
- [ ] **Implements Attribute Refreshes During Read**: All attributes available in the API should have [`d.Set()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.Set) called their values in the Terraform state during the `Read` function.
- [ ] **Implements Error Checks with Non-Primitive Attribute Refreshes**: When using [`d.Set()`](https://godoc.org/github.com/hashicorp/terraform/helper/schema#ResourceData.Set) with non-primitive types (`schema.TypeList`, `schema.TypeSet`, or `schema.TypeMap`), perform error checking to [prevent issues where the code is not properly able to refresh the Terraform state](https://www.terraform.io/docs/extend/best-practices/detecting-drift.html#error-checking-aggregate-types).
- [ ] **Implements Import Acceptance Testing and Documentation**: Support for resource import (`Importer` in resource schema) must include `ImportState` acceptance testing (see also the [Acceptance Testing Guidelines](#acceptance-testing-guidelines) below) and `## Import` section in resource documentation.
- [ ] **Implements Customizable Timeouts Documentation**: Support for customizable timeouts (`Timeouts` in resource schema) must include `## Timeouts` section in resource documentation.
- [ ] **Implements State Migration When Adding New Virtual Attribute**: For new "virtual" attributes (those only in Terraform and not in the API), the schema should implement [State Migration](https://www.terraform.io/docs/extend/resources.html#state-migrations) to prevent differences for existing configurations that upgrade.
- [ ] **Uses Okta Go SDK Types**: Use available SDK structs instead of implementing custom types with indirection.
- [ ] **Uses Existing Validation Functions**: Schema definitions including `ValidateFunc` for attribute validation should use available [Terraform `helper/validation` package](https://godoc.org/github.com/hashicorp/terraform/helper/validation) functions. `All()`/`Any()` can be used for combining multiple validation function behaviors.
- [ ] **Skips Exists Function**: Implementing a resource `Exists` function is extraneous as it often duplicates resource `Read` functionality. Ensure `d.SetId("")` is used to appropriately trigger resource recreation in the resource `Read` function.
- [ ] **Skips id Attribute**: The `id` attribute is implicit for all Terraform resources and does not need to be defined in the schema.

The below are style-based items that _may_ be noted during review and are recommended for simplicity, consistency, and quality assurance:

- [ ] **Avoids CustomizeDiff**: Usage of `CustomizeDiff` is generally discouraged.
- [ ] **Implements Error Message Context**: Returning errors from resource `Create`, `Read`, `Update`, and `Delete` functions should include additional messaging about the location or cause of the error for operators and code maintainers by wrapping with [`fmt.Errorf()`](https://godoc.org/golang.org/x/exp/errors/fmt#Errorf).
  - An example `Delete` API error: `return fmt.Errorf("error deleting {SERVICE} {THING} (%s): %s", d.Id(), err)`
  - An example `d.Set()` error: `return fmt.Errorf("error setting {ATTRIBUTE}: %s", err)`
- [ ] **Uses Elem with TypeMap**: While provider schema validation does not error when the `Elem` configuration is not present with `Type: schema.TypeMap` attributes, including the explicit `Elem: &schema.Schema{Type: schema.TypeString}` is recommended.
- [ ] **Uses American English for Attribute Naming**: For any ambiguity with attribute naming, prefer American English over British English. e.g. `color` instead of `colour`.
- [ ] **Skips Timestamp Attributes**: Generally, creation and modification dates from the API should be omitted from the schema.
- [ ] **Skips Error() Call with Okta Go SDK Error Objects**: Error objects do not need to have `Error()` called.

#### Acceptance Testing Guidelines

The below are required items that will be noted during submission review and prevent immediate merging:

- [ ] **Implements CheckDestroy**: Resource testing should include a `CheckDestroy` function (typically named `testAccCheckAws{SERVICE}{RESOURCE}Destroy`) that calls the API to verify that the Terraform resource has been deleted or disassociated as appropriate. More information about `CheckDestroy` functions can be found in the [Extending Terraform TestCase documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html#checkdestroy).
- [ ] **Implements Exists Check Function**: Resource testing should include a `TestCheckFunc` function (typically named `testAccCheckAws{SERVICE}{RESOURCE}Exists`) that calls the API to verify that the Terraform resource has been created or associated as appropriate. Preferably, this function will also accept a pointer to an API object representing the Terraform resource from the API response that can be set for potential usage in later `TestCheckFunc`. More information about these functions can be found in the [Extending Terraform Custom Check Functions documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html#checkdestroy).
- [ ] **Excludes Provider Declarations**: Test configurations should not include `provider "okta" {...}` declarations. If necessary, only the provider declarations in `provider_test.go` should be used for multiple account/region or otherwise specialized testing.
- [ ] **Passes in us-west-2 Region**: Tests default to running in `us-west-2` and at a minimum should pass in that region or include necessary `PreCheck` functions to skip the test when ran outside an expected environment.
- [ ] **Uses resource.ParallelTest**: Tests should utilize [`resource.ParallelTest()`](https://godoc.org/github.com/hashicorp/terraform/helper/resource#ParallelTest) instead of [`resource.Test()`](https://godoc.org/github.com/hashicorp/terraform/helper/resource#Test) except where serialized testing is absolutely required.
- [ ] **Uses fmt.Sprintf()**: Test configurations preferably should to be separated into their own functions (typically named `testAccAws{SERVICE}{RESOURCE}Config{PURPOSE}`) that call [`fmt.Sprintf()`](https://golang.org/pkg/fmt/#Sprintf) for variable injection or a string `const` for completely static configurations. Test configurations should avoid `var` or other variable injection functionality such as [`text/template`](https://golang.org/pkg/text/template/).
- [ ] **Uses Randomized Infrastructure Naming**: Test configurations that utilize resources where a unique name is required should generate a random name. Typically, this is created via `rName := acctest.RandomWithPrefix("tf-acc-test")` in the acceptance test function before generating the configuration.

For resources that support import, the additional item below is required that will be noted during submission review and prevent immediate merging:

- [ ] **Implements ImportState Testing**: Tests should include an additional `TestStep` configuration that verifies resource import via `ImportState: true` and `ImportStateVerify: true`. This `TestStep` should be added to all possible tests for the resource to ensure that all infrastructure configurations are properly imported into Terraform.

The below are style-based items that _may_ be noted during review and are recommended for simplicity, consistency, and quality assurance:

- [ ] **Uses Builtin Check Functions**: Tests should utilize already available check functions, e.g. `resource.TestCheckResourceAttr()`, to verify values in the Terraform state over creating custom `TestCheckFunc`. More information about these functions can be found in the [Extending Terraform Builtin Check Functions documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/teststep.html#builtin-check-functions).
- [ ] **Uses TestCheckResourceAttrPair() for Data Sources**: Tests should utilize [`resource.TestCheckResourceAttrPair()`](https://godoc.org/github.com/hashicorp/terraform/helper/resource#TestCheckResourceAttrPair) to verify values in the Terraform state for data sources attributes to compare them with their expected resource attributes.
- [ ] **Excludes Timeouts Configurations**: Test configurations should not include `timeouts {...}` configuration blocks except for explicit testing of customizable timeouts (typically very short timeouts with `ExpectError`).
- [ ] **Implements Default and Zero Value Validation**: The basic test for a resource (typically named `TestAccAws{SERVICE}{RESOURCE}_basic`) should utilize available check functions, e.g. `resource.TestCheckResourceAttr()`, to verify default and zero values in the Terraform state for all attributes. Empty/missing configuration blocks can be verified with `resource.TestCheckResourceAttr(resourceName, "{ATTRIBUTE}.#", "0")` and empty maps with `resource.TestCheckResourceAttr(resourceName, "{ATTRIBUTE}.%", "0")`

The below are location-based items that _may_ be noted during review and are recommended for consistency with testing flexibility. Resource testing is expected to pass across multiple Okta environments supported by the Terraform Okta Provider (e.g. Okta Standard and Okta GovCloud (US)). Contributors are not expected or required to perform testing outside of Okta Standard, e.g. running only in the `us-west-2` region is perfectly acceptable, however these are provided for reference:

### Writing Acceptance Tests

Terraform includes an acceptance test harness that does most of the repetitive
work involved in testing a resource. For additional information about testing
Terraform Providers, see the [Extending Terraform documentation](https://www.terraform.io/docs/extend/testing/index.html).

#### Acceptance Tests Often Cost Money to Run

Because acceptance tests create real resources, they often cost money to run.
Because the resources only exist for a short period of time, the total amount
of money required is usually a relatively small. Nevertheless, we don't want
financial limitations to be a barrier to contribution, so if you are unable to
pay to run acceptance tests for your contribution, mention this in your
pull request. We will happily accept "best effort" implementations of
acceptance tests and run them for you on our side. This might mean that your PR
takes a bit longer to merge, but it most definitely is not a blocker for
contributions.

#### Acceptance Tests With VCR

The provider's acceptance tests are integrated with
[go-vcr](https://github.com/dnaeon/go-vcr) to allow running, fast, valid, off
the wire acceptance tests. The signal for VCR mode is the ENV var
`OKTA_VCR_TF_ACC` with a value of `record` or `play`. If the mode is `play`
then the named test is run for each of the cassettes (YAML files) in the
`test/fixtures/vcr/{test-name}/` directory. If the mode is `record` then the
ENV var `OKTA_VCR_CASSETTE` with a value for the cassette name must be set. All
of the HTTP requests/responses are recorded to that cassette file. Finally, if
the `test/fixtures/vcr/{test-name}/` directory or
`test/fixtures/vcr/{test-name}/{cassete-name}`file does not exist then that
named test is marked a skip.

The harness will not overwrite an existing cassette file when it is in record
mode.  The developer must delete that cassette from the file system first
before it is re-recorded.

An important subtlety to be called out here is that cassettes can be recorded
for different orgs. Therefore using an intelligent naming convention for
`OKTA_VCR_CASSETTE` new cassettes can be recorded an org with specific feature
flags enabled on it. Thereafter, that org doesn't need to be live any further
to have its tests run as they are replayed with VCR. If you add new cassettes
to the repo please add notes about them in the `test/fixtures/vcr/CASSETTES.md`
file.

Examples:

Run all VCR enabled tests
```
OKTA_VCR_TF_ACC=play make testacc
# or
make test-play-vcr-acc
```

Run a single VCR enabled test
```
OKTA_VCR_TF_ACC=play make testacc TEST=./okta TESTARGS='-run=TestAccOktaGroup_crud'
# or
make test-play-vcr-acc TEST=./okta TESTARGS='-run=TestAccOktaGroup_crud'
```

Record a new VCR cassette for specific test
```
OKTA_VCR_CASSETTE=oie-with-feature-x OKTA_VCR_TF_ACC=record make testacc TEST=./okta TESTARGS='-run=TestAccOktaGroup_crud'
# or
OKTA_VCR_CASSETTE=oie-with-feature-x make test-record-vcr-acc TEST=./okta TESTARGS='-run=TestAccOktaGroup_crud'
# NOTE cassette file dropped in:
# test/fixtures/vcr/TestAccOktaGroup_crud/oie-with-feature-x.yaml
```

Record a new VCR cassette for all tests
```
OKTA_VCR_CASSETTE=oie-with-feature-x OKTA_VCR_TF_ACC=record make testacc
# or
OKTA_VCR_CASSETTE=oie-with-feature-x make test-record-vcr-acc
```

#### Running an Acceptance Test

Acceptance tests can be run using the `testacc` target in the Terraform
`Makefile`. The individual tests to run can be controlled using a regular
expression. Prior to running the tests, provider configuration details such as
access keys must be made available as environment variables.

For example, to run an acceptance test against the Okta
provider, the following environment variables must be set:

```sh
export OKTA_API_TOKEN=...
export OKTA_ORG_NAME=...
export OKTA_BASE_URL=...
```

Tests can then be run by specifying the target provider, and a regular expression defining the tests to run:

```sh
$ make testacc TEST=./okta TESTARGS='-run=TestAccOktaOAuthApp_crud'
==> Checking that code complies with gofmt requirements...
TF_ACC=1 go test ./okta -v -run=TestAccOktaOAuthApp_crud -timeout 120m
=== RUN   TestAccOktaOAuthApp_crud
--- PASS: TestAccOktaOAuthApp_crud (26.56s)
PASS
ok  	github.com/okta/terraform-provider-okta/okta	26.607s
```

Entire resource test suites can be targeted by using the naming convention to
write the regular expression. For example, to run all tests of the
`okta_user_schema` resource rather than just the update test, you can start
testing like this:

```sh
$ make testacc TEST=./okta TESTARGS='-run=TestAccOktaUserSchema'
==> Checking that code complies with gofmt requirements...
TF_ACC=1 go test ./okta -v -run=TestAccOktaUserSchema -timeout 120m
=== RUN   TestAccOktaUserSchema_crud
--- PASS: TestAccOktaUserSchema_crud (15.06s)
=== RUN   TestAccOktaUserSchema_arrayString
--- PASS: TestAccOktaUserSchema_arrayString (12.70s)
PASS
ok  	github.com/okta/terraform-provider-okta/okta	55.619s
```

#### Forced clean out of dangling acceptance test resources

Sometimes the acceptance testing framework can exit leaving dangling resources.
Run the forced sweeper test to clean them all out.

```
OKTA_ACC_TEST_FORCE_SWEEPERS=1 TF_LOG=warn make testacc TEST=./okta TESTARGS='-run=TestRunForcedSweeper'
```

#### Writing an Acceptance Test

Terraform has a framework for writing acceptance tests which minimises the
amount of boilerplate code necessary to use common testing patterns. The entry
point to the framework is the `resource.ParallelTest()` function.

Tests are divided into `TestStep`s. Each `TestStep` proceeds by applying some
Terraform configuration using the provider under test, and then verifying that
results are as expected by making assertions using the provider API. It is
common for a single test function to exercise both the creation of and updates
to a single resource. Most tests follow a similar structure.

1. Pre-flight checks are made to ensure that sufficient provider configuration
   is available to be able to proceed - for example in an acceptance test
   targeting Okta, `OKTA_BASE_URL`, `OKTA_API_TOKEN`, and `OKTA_ORG_NAME` must be set prior
   to running acceptance tests. This is common to all tests exercising a single
   provider.

Each `TestStep` is defined in the call to `resource.ParallelTest()`. Most assertion
functions are defined out of band with the tests. This keeps the tests
readable, and allows a reuse of assertion functions across different tests of the
same type of resource. The definition of a complete test looks like this:

```go
func TestAccOktaGroups_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", oktaGroup)
	mgr := newFixtureManager("okta_group")
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)
	addUsersConfig := mgr.GetFixtures("okta_group_with_users.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oktaGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAccDifferent")),
			},
			{
				Config: addUsersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "4"),
				),
			},
		},
	})
}
```

When executing the test, the following steps are taken for each `TestStep`:

2. The Terraform configuration required for the test is applied. This is
   responsible for configuring the resource under test, and any dependencies it
   may have. For example, to test the `okta_group` resource, a valid configuration with the requisite fields is required. This results in configuration which looks like this:

   ```hcl
   resource okta_group test {
       name        = "testAcc_replace_with_uuid"
       description = "testing"
   }
   ```

3. Assertions are run using the provider API. These use the provider API
   directly rather than asserting against the resource state. For example, to
   verify that the `okta_group` described above was created
   successfully, a test function like this is used:

   ```go
   func doesGroupExist(id string) (bool, error) {
       client := getOktaClientFromMetadata(testAccProvider.Meta())
       _, response, err := client.Group.GetGroup(id, nil)

       return doesResourceExist(response, err)
   }
   ```

1. The resources created by the test are destroyed. This step happens
   automatically, and is the equivalent of calling `terraform destroy`.

1. Assertions are made against the provider API to verify that the resources
   have indeed been removed. If these checks fail, the test fails and reports
   "dangling resources". The code to ensure that the `okta_user` shown
   above has been destroyed looks like this:

   ```go
   CheckDestroy: createCheckResourceDestroy(oktaGroup, doesGroupExist)
   ```

[website]: https://github.com/okta/terraform-provider-okta/tree/main/website
[acctests]: https://github.com/hashicorp/terraform#acceptance-tests
[ml]: https://groups.google.com/group/terraform-tool
