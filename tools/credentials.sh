#!/usr/bin/env bash
PHIUSER=${PHIUSER:-someusername}
PHIPASS=${PHIPASS:-somesupersecretpassword}

if [[ -z PHITOKEN ]]; then
    $PHITOKEN=$(./login | tail -1 | jq -r .token)
    export PHITOKEN="${PHITOKEN}"
fi
