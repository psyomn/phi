#!/usr/bin/env bash
source "credentials.sh"
curl -s -i -X POST "localhost:9876/register" -d "{\"username\":\"$PHIUSER\",\"password\":\"$PHIPASS\"}"
