#!/usr/bin/env bash
source "credentials.sh"

curl -X POST \
     -H "Authorization: Token ${PHITOKEN}" \
     --data-binary @tux.png \
     localhost:9876/upload
