# Maintaining the Terraform Okta Provider

<!-- TOC depthFrom:2 -->

- [Pull Requests](#pull-requests)
    - [Pull Request Review Process](#pull-request-review-process)
        - [Dependency Updates](#dependency-updates)
            - [Okta Go SDK Updates](#okta-go-sdk-updates)
            - [golangci-lint Updates](#golangci-lint-updates)
            - [Terraform Plugin SDK Updates](#terraform-plugin-sdk-updates)
            - [tfproviderdocs Updates](#tfproviderdocs-updates)
            - [tfproviderlint Updates](#tfproviderlint-updates)
    - [Pull Request Merge Process](#pull-request-merge-process)
    - [Pull Request Types to CHANGELOG](#pull-request-types-to-changelog)
- [Release Process](#release-process)

<!-- /TOC -->

## Pull Requests

### Pull Request Review Process

Notes for each type of pull request are (or will be) available in subsections below.

- If you plan to be responsible for the pull request through the merge/closure process, assign it to yourself
- Add `bug`, `enhancement`, `new-data-source`, `new-resource`, or `technical-debt` labels to match expectations from change
- Perform a quick scan of open issues and ensure they are referenced in the pull request description (e.g. `Closes #1234`, `Relates #5678`). Edit the description yourself and mention this to the author:

```markdown
This pull request appears to be related to/solve #1234, so I have edited the pull request description to denote the issue reference.
```

- Review the contents of the pull request and ensure the change follows the relevant section of the [Contributing Guide](https://github.com/okta/terraform-provider-okta/blob/main/.github/CONTRIBUTING.md#checklists-for-contribution)
- If the change is not acceptable, leave a long form comment about the reasoning and close the pull request
- If the change is acceptable with modifications, leave a pull request review marked using the `Request Changes` option (for maintainer pull requests with minor modification requests, giving feedback with the `Approve` option is recommended, so they do not need to wait for another round of review)
- If the author is unresponsive for changes (by default we give two weeks), determine importance and level of effort to finish the pull request yourself including their commits or close the pull request
- Run relevant acceptance testing ([locally](https://github.com/okta/terraform-provider-okta/blob/main/.github/CONTRIBUTING.md#running-an-acceptance-test) or in TeamCity) against an Okta account to ensure no new failures are being introduced
- Approve the pull request with a comment outlining what steps you took that ensure the change is acceptable, e.g. acceptance testing output

``````markdown
Looks good, thanks @username! :rocket:

Output from acceptance testing in Okta:

```
--- PASS: TestAcc...
--- PASS: TestAcc...
```
``````

#### Dependency Updates

##### Okta Go SDK Updates

Almost exclusively, `github.com/okta/okta-sdk-go` updates are additive in nature. It is generally safe to only scan through them before approving and merging. If you have any concerns about any of the service client updates such as suspicious code removals in the update, or deprecations introduced, run the acceptance testing for potentially affected resources before merging.

##### Terraform Plugin SDK Updates

Except for trivial changes, run the full acceptance testing suite against the pull request and verify there are no new or unexpected failures.

##### tfproviderdocs Updates

Merge if CI passes.

##### tfproviderlint Updates

Merge if CI passes.

### Pull Request Merge Process

- Add this pull request to the upcoming release milestone
- Add any linked issues that will be closed by the pull request to the same upcoming release milestone
- Merge the pull request
- Delete the branch (if the branch is on this repository)
- Determine if the pull request should have a CHANGELOG entry by reviewing the [Pull Request Types to CHANGELOG section](#pull-request-types-to-changelog). If so, update the repository `CHANGELOG.md` by directly committing to the `main` branch (e.g. editing the file in the GitHub web interface). See also the [Extending Terraform documentation](https://www.terraform.io/docs/extend/best-practices/versioning.html) for more information about the expected CHANGELOG format.
- Leave a comment on any issues closed by the pull request noting that it has been merged and when to expect the release containing it, e.g.

```markdown
The fix for this has been merged and will release with a version X.Y.Z of the Terraform Okta Provider, expected in the XXX timeframe.
```

### Pull Request Types to CHANGELOG

The CHANGELOG is intended to show operator-impacting changes to the codebase for a particular version. If every change or commit to the code resulted in an entry, the CHANGELOG would become less useful for operators. The lists below are general guidelines on when a decision needs to be made to decide whether a change should have an entry.

Changes that should have a CHANGELOG entry:

- New Resources and Data Sources
- New full-length documentation guides (e.g. EKS Getting Started Guide, IAM Policy Documents with Terraform)
- Resource and provider bug fixes
- Resource and provider enhancements
- Deprecations
- Removals

Changes that may have a CHANGELOG entry:

- Dependency updates: If the update contains relevant bug fixes or enhancements that affect operators, those should be called out.

Changes that should _not_ have a CHANGELOG entry:

- Resource and provider documentation updates
- Testing updates

## Release Process

- Create a milestone for the next release after this release (generally, the next milestone will be a minor version increase unless previously decided for a major or patch version)
- Check the existing release milestone for open items and either work through them or move them to the next milestone
- Run the HashiCorp (non-OSS) TeamCity release job with the `DEPLOYMENT_TARGET_VERSION` matching the expected release milestone and `DEPLOYMENT_NEXT_VERSION` matching the next release milestone
- Wait for the TeamCity release job and CircleCI website deployment jobs to complete either by watching the build logs or Slack notifications
- Close the release milestone
- Create a new GitHub release with the release title exactly matching the tag and milestone (e.g. `v2.22.0`) and copy the notes from the CHANGELOG to the release notes. This will trigger [HashiBot](https://github.com/apps/hashibot) release comments.
