#!/usr/bin/env bash
set -eux

function solc() {
	docker run -u 1000:1000 -v $(pwd):/sources -w /sources ethereum/solc:0.5.12 "$@"
}

NAME="osm"
PACKAGE="maker/osm"
CONTRACT="OSM"

mkdir -p "./contracts/${PACKAGE}/internal"

go run ./flatcat ~/maker/osm/src/osm.sol > "./contracts/${PACKAGE}/internal/${NAME}.sol"

solc "./contracts/${PACKAGE}/internal/${NAME}.sol" --output-dir "./contracts/${PACKAGE}/internal" --abi --bin

mv "./contracts/${PACKAGE}/internal/${CONTRACT}.abi" "./contracts/${PACKAGE}/internal/${NAME}.abi"
mv "./contracts/${PACKAGE}/internal/${CONTRACT}.bin" "./contracts/${PACKAGE}/internal/${NAME}.bin"

shopt -s extglob
rm ./catalog/${PACKAGE}/internal/!(${NAME}).*

abigen --abi "./contracts/${PACKAGE}/internal/${NAME}.abi" --pkg internal --type "${NAME}" --out "./contracts/${PACKAGE}/internal/${NAME}.go"