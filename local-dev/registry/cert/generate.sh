#!/usr/bin/env bash

set -eux

openssl req -new -newkey rsa:4096 -x509 -sha256 -days 3650 -nodes -out domain.crt -keyout domain.key \
 -addext "subjectAltName = DNS:registry.local.tokyotech.org"
