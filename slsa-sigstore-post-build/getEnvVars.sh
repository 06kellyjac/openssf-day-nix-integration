#!/usr/bin/env sh

GIVEN_PATH=$1

DRV="$(nix derivation show "$GIVEN_PATH" | jq 'keys | .[0]' -r)"
OUT="$(nix derivation show "$GIVEN_PATH" | jq 'first(.[]) | .outputs | map(.path) | join(" ")' -r)"

echo "DRV_PATH=\"$DRV\" OUT_PATHS=\"$OUT\""
