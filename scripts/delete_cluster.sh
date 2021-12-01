#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

terraform destroy --auto-approve

rm cluster.tf

if [ -d ".terraform" ]
then
    rm -rf .terraform
fi

if [ -d "terraform.tfstate.d" ]
then
    rm -rf terraform.tfstate.d
fi

