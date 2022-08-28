#!/usr/bin/env bash

set -e
SCRIPT=$(realpath "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

echo "Setting development override for okta/okta provider"
export TF_CLI_CONFIG_FILE="$SCRIPTPATH/dev.tfrc"
