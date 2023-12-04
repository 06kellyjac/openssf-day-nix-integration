#!/usr/bin/env sh

rekor-cli get --rekor_server http://localhost:3000 --log-index "$1" --format json | jq | bat -pl json
