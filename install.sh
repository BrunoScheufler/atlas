#!/usr/bin/env bash

set -eo pipefail

go mod download
go build -o atlas ./cli
mv atlas /usr/local/bin/atlas
