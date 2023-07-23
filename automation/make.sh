#!/usr/bin/env bash

set -e

ARGCOUNT=$#

EXEC_PATH=$(dirname "$(realpath "$0")")
PROJECT_PATH="$(dirname $EXEC_PATH)"

CONTAINER_WORKSPACE="/workspace/template-cni"

: "${CONTAINER_CMD:="podman"}"
: "${CONTAINER_IMG:=" quay.io/projectquay/golang:1.19"}"

: "${DISABLE_IPV6_IN_CONTAINER:=0}"

test -t 1 && USE_TTY="-t"

options=$(getopt --options "" \
    --long lint,build,fmt,unit-test,help\
    -- "${@}")
eval set -- "$options"
while true; do
    case "$1" in
    --lint)
        OPT_LINT=1
        ;;
    --build)
        OPT_BUILD=1
        ;;
    --fmt)
        OPT_FMT=1
        ;;
    --unit-test)
        OPT_UTEST=1
        ;;
    --help)
        set +x
        echo "$0 [--lint] [--build] [--fmt] [--unit-test]"
        exit
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

if [ "${ARGCOUNT}" -eq "0" ]; then
    OPT_BUILD=1
    OPT_FMT=1
    OPT_LINT=1
    OPT_UTEST=1
fi

if [ -n "${OPT_LINT}" ]; then
    golangci_lint_version=v1.53.3
    if [ ! -f "$(go env GOPATH)"/bin/golangci-lint ]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)"/bin $golangci_lint_version
    fi
    golangci-lint run -v cmd/... pkg/...
fi

if [ -n "${OPT_BUILD}" ]; then
    go build -o ./bin/plugin -v ./cmd/...
fi

if [ -n "${OPT_FMT}" ]; then
        unformatted=$(gofmt -l ./cmd ./pkg)
        test -z "$unformatted" || (echo "Unformatted: $unformatted" && false)
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./cmd/... ./pkg/...
fi
