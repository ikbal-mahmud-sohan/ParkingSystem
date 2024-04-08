#!/usr/bin/env bash

set -euxo pipefail

source ./devenv.sh

GOBIN=$(pwd)/bin go install \
        github.com/gunk/gunk \

$(pwd)/bin/gunk format ./...

$(pwd)/bin/gunk generate ./...

if [[ -d swagger.json ]]; then rm -r rm swagger.json; fi

swagger-combine config.json -o swagger.json