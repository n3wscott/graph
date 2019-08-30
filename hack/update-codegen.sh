#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CODEGEN_PKG=./vendor/k8s.io/code-generator

# Only deepcopy the Duck types, as they are not real resources.
${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
  github.com/n3wscott/graph/pkg/client github.com/n3wscott/graph/pkg/apis \
  "duck:v1alpha1"

# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh