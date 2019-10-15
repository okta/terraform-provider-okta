# Steps to Release

We're following [semantic versioning](https://semver.org/) approach.

## Create a Release PR

Assemble all the meaningful changes since the last release into the [CHANGELOG.md](CHANGELOG.md) file.

## Merge Release PR

Verify that the acceptance test suite has passed for the release PR, then merge the PR.

## Tag the release

```
git tag -am "terraform-provider-okta_vX.Y.Z_x5 release" vX.Y.Z
git push --tags
```

## Goreleaser

We run [goreleaser](https://goreleaser.com/) via travis, this should run against master on commit. Here is how to run it manually:

```
GITHUB_TOKEN=xxx goreleaser release --rm-dist
```

## See the release in Github

You can find the releases in Github, e.g. [v3.0.0](https://github.com/terraform-providers/terraform-provider-okta/releases/tag/v3.0.0).
