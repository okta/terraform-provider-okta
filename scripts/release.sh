#!/usr/bin/env bash
set -e

lastTag=$(git tag --sort=v:refname | tail -1)
patch="${lastTag##*.}"
versionTag="${lastTag%.*}.$((patch+1))"
git tag "$versionTag"
git push --tags
versionTag="${versionTag#*v}"
goreleaser release
TERRAFORM_PLUGINS_PROFILE=prod TERRAFORM_PLUGINS_BUCKET=articulate-terraform-providers VERSION="$versionTag" make ship
