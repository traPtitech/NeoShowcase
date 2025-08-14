#!/usr/bin/env bash

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <chart name> <destination>"
  exit 1
fi

set -euxo pipefail

CHART=$1
DEST=$2

helm package "$CHART"
helm push "$CHART"-*.tgz "$DEST"
rm "$CHART"-*.tgz
