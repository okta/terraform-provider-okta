#!/bin/zsh

RESOURCE_ID="${1:-iam7ghmsqfSEtsDmK1d7}"

if [ ! -f terraform.auto.tfvars ]; then
    cat << EOD
Add a terraform.auto.tfvars with the following data
partition = ""
org_name  = ""
org_id    = ""
api_token = ""
EOD
    exit 1
fi

cd `dirname $0`/../../
make

cd -
rm -rf .terraform* terraform.tfstate*

terraform init
terraform import okta_resource_set.my_resources ${RESOURCE_ID}
terraform show -json terraform.tfstate > state.json

env TF_LOG=INFO terraform apply