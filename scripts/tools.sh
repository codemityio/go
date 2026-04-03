#!/bin/bash

set -e

source scripts/func.sh

# Note: the set [-/+] x is purely there to turn on and off outputting of the commands being executed.
if [ "${DEBUG}" = "true" ]; then
  set -x
fi

mkdir -p tmp var

case "$1" in

"prep")
  go install github.com/"${VENDOR}"/auxilium@latest
  go install github.com/"${VENDOR}"/notatio@latest
  go install github.com/"${VENDOR}"/goforma@latest
  ;;

"cmd")
  docker run --rm \
    --name "${BASE_NAME}-cmd" \
    -e DEBUG \
    -e GOBIN="${PWD}/bin" \
    -e CGO_ENABLED=1 \
    -v "${PWD}:${PWD}" \
    -w "${PWD}" \
    "${VENDOR}"/golang-dev:latest sh -c "${COMMAND}"
  ;;

"cleanup")
  find . -iname '*_mock.go' -exec rm {} \;
  find . -iname '*.so' -exec rm {} \;
  git clean -dXf pkg internal cmd
  rm -Rf tmp var bin
  ;;

"update")
  COMMAND="go mod init github.com/${VENDOR}/$(basename "${PWD}")"
  export COMMAND
  rm -rf go.* && scripts/tools.sh cmd
  scripts/tools.sh fmt
  ;;

"gen")
  go generate -v ./...
  ;;

"fmt")
  go mod tidy
  goimports -w .
  gofumpt -l -w .
  go vet ./...
  ;;

"statan")
  statan "./..." ""
  ;;

"statan-fix")
  statan "./..." "--fix"
  ;;

"test")
  go test -failfast -v -covermode=count -coverprofile=tmp/coverage.out ./...
  ;;

"test-race")
  go test -race -failfast -v -covermode=atomic -coverprofile=tmp/coverage.out ./...
  ;;

"cov")
  go tool cover -func="tmp/coverage.out" -o tmp/coverage.in
  goforma badge \
    --document-path=README.md \
    --id=coverage-badge-do-not-edit \
    coverage \
    --cov-file-path=tmp/coverage.in \
    --minimum="${MINIMUM_COVERAGE}"
  ;;

"cov-report")
  go tool cover -func=tmp/coverage.out
  ;;

"cov-open")
  go tool cover -html=tmp/coverage.out -o tmp/coverage.html
  open tmp/coverage.html
  ;;

"diff")
  (git diff --quiet && git diff --cached --quiet && [ -z "$(git ls-files --others --exclude-standard)" ]) || {
    echo "error: changes detected..."
    exit 1
  }
  ;;

"version")
  tag=$(git tag -l | sort -V | tail -n1)
  echo "${tag:-latest}"
  ;;

"next")
  # get latest v-tag (supports vX.Y.Z and vX.Y.Z-preN) and bump prerelease or patch
  latest=$(git tag -l "v*" | sort -V | tail -n1)

  if [[ $latest =~ ^v([0-9]+\.[0-9]+\.[0-9]+)-([a-zA-Z]+)([0-9]+)$ ]]; then
    next="v${BASH_REMATCH[1]}-${BASH_REMATCH[2]}$((BASH_REMATCH[3] + 1))"
  elif [[ $latest =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    next="v${BASH_REMATCH[1]}.${BASH_REMATCH[2]}.$((BASH_REMATCH[3] + 1))"
  else
    echo "error: unsupported tag format: $latest"
    exit 1
  fi

  git tag "$next"
  echo "Tagged: $next"
  ;;

*)
  echo "error: incorrect '$1' command..."
  ;;

esac
