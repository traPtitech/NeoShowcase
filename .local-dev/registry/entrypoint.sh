#!/bin/sh

set -e

CONFIG="$1"
${GC_SCRIPT:-/gc.sh} "$CONFIG" &

set -- registry serve "$@"
exec "$@"
