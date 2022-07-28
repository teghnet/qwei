#!/usr/bin/env bash
set -xeuo pipefail
_here="$(cd "${0%/*}" && pwd)"
cd "$_here"

function solc() {
	docker run -u 1000:1000 -v $(pwd):/sources -w /sources ethereum/solc:0.6.12 "$@"
}

NAME="changelog"
CONTRACT="ChainLog"

mkdir -p "./internal"
rm -f ./internal/*

go run ../../../flatcat ~/maker/dss-chain-log/src/ChainLog.sol > "$_here/${NAME}.sol"

solc "./${NAME}.sol" --output-dir "./internal" --abi --bin --optimize --optimize-runs=200

jq < "./internal/${CONTRACT}.abi" > "./${NAME}.abi"
rm "./internal/${CONTRACT}.abi"
mv "./internal/${CONTRACT}.bin" "./${NAME}.bin"

abigen --abi "./${NAME}.abi" --pkg internal --type "${NAME}" --out "./internal/${NAME}.go"
