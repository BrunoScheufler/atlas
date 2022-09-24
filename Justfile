set positional-arguments

default:
    just --list

dev:
  #!/usr/bin/env bash

  set -eo pipefail

  go mod download
  go build --ldflags "-X main.version=dev" -o atlas ./cli
  rm -rf /usr/local/bin/atlas
  mv atlas /usr/local/bin/atlas
