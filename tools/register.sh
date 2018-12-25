#!/usr/bin/env bash
USER=${USER:-bob}
PASS=${PASS:-secret}
curl -s -i -X POST "localhost:9876/register" -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}"
