#!/usr/bin/env bash

set -eux

TMP_DIR=$(mktemp -d)
trap 'sudo rm -r "$TMP_DIR"' EXIT
cd "$TMP_DIR"

case $(uname -s) in
  "Linux") OS="linux";;
  "Darwin") OS="osx";;
  *) echo "unknown os" && exit 1;;
esac

ARCH=$(uname -m)

if [ $ARCH = "arm64" ]; then
  ARCH="aarch_64"
fi

curl -sSLOf https://github.com/protocolbuffers/protobuf/releases/download/v"$PROTOC_VERSION"/protoc-"$PROTOC_VERSION"-"$OS"-"$ARCH".zip
unzip protoc-"$PROTOC_VERSION"-"$OS"-"$ARCH".zip

sudo install bin/protoc /usr/local/bin/
