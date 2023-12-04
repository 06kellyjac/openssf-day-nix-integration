#! /usr/bin/env nix-shell
#! nix-shell -i bash -p bash skopeo cosign rekor-cli

./slsa-sigstore-post-build
