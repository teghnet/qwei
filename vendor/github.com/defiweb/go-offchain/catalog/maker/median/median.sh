#!/usr/bin/env bash
set -eux

function solc() {
	docker run -u 1000:1000 -v $(pwd):/sources -w /sources ethereum/solc:0.5.12 "$@"
}

NAME="median"
PACKAGE="maker/median"
CONTRACT="MedianWSTETHUSD"

mkdir -p "./contracts/${PACKAGE}/internal"

solc "./contracts/${PACKAGE}/internal/${NAME}.sol" --output-dir "./contracts/${PACKAGE}/internal" --abi --bin --overwrite

mv "./contracts/${PACKAGE}/internal/${CONTRACT}.abi" "./contracts/${PACKAGE}/internal/${NAME}.abi"
mv "./contracts/${PACKAGE}/internal/${CONTRACT}.bin" "./contracts/${PACKAGE}/internal/${NAME}.bin"

shopt -s extglob
rm -f ./catalog/${PACKAGE}/internal/!(${NAME}).*

abigen --abi "./contracts/${PACKAGE}/internal/${NAME}.abi" --pkg internal --type "${NAME}" --out "./contracts/${PACKAGE}/internal/${NAME}.go"