#!/bin/sh

while true; do
  sleep 3600
  echo "[gc.sh] Running garbage collection..."
  registry garbage-collect "$1" --delete-untagged
  echo "[gc.sh] Garbage collection complete."
done
