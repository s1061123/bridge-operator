#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

bash vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/s1061123/bridge-operator/pkg/client github.com/s1061123/bridge-operator/pkg/apis \
  openshift.io:v1alpha1 \
  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt
