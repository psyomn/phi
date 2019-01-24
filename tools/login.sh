#!/usr/bin/env bash
USER=${USER:-bob}
PASS=${PASS:-pass}
curl -s -i -X POST "localhost:9876/login" -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}"
