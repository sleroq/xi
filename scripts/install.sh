#!/bin/bash

set -e

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

cd "$(dirname "$SCRIPTPATH")"

BIN="$(go env GOBIN)"

if [ -z "$BIN" ]
then
    BIN="$(go env GOPATH)"/bin
fi

if [ -z "$BIN" ]
then
    echo "GOBIN and GOPATH are not set.
Add binary (built/xi) to your \$PATH manually."
    exit 1
fi

scripts/build.sh

cp built/xi "$BIN"
