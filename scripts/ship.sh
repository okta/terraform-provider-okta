# Articulate Specific ship pipeline
#!/usr/bin/env bash

version=$(git tag --sort=v:refname | tail -1)

mkdir -p ~/.terraform.d/plugins/

gox -osarch="linux/amd64 darwin/amd64 windows/amd64" \
  -output="${HOME}/.terraform.d/plugins/{{.OS}}_{{.Arch}}/terraform-provider-okta_${version}_x4" .

echo "copying terraform-provider-okta_${version}_x4 to s3://${TERRAFORM_PLUGINS_BUCKET}"
aws s3 cp \
  "${HOME}/.terraform.d/plugins/linux_amd64/terraform-provider-okta_${version}_x4" \
  "s3://${TERRAFORM_PLUGINS_BUCKET}" \
  --profile "${TERRAFORM_PLUGINS_PROFILE}"
