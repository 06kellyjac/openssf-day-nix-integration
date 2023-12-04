# Cheatsheet

```sh
# testing
DRV_PATH="/nix/store/zjfzky5nlq7cn9fjwvivvaj0snasfiwn-basic-c.drv" OUT_PATHS="/nix/store/8ar96fb4lhsmz42bpgj7czja35sp6qya-basic-c /nix/store/9abv367ibw74q4dmqcdcq01x3j4xvp20-basic-c-extras /nix/store/a2gn1mh340k3773ibaf4q8j5aazaqbxj-basic-c-lib" go run . | jq

# container tar.gz
DRV_PATH="/nix/store/gspmd3alks44drwfkx75pm8h381lazmd-hello.tar.gz.drv" OUT_PATHS="/nix/store/29cxziviylsf2rhpijrz9sizl53n9h3s-hello.tar.gz" go run . | jq


./getEnvVars.sh ./path
```

```
# FOD
/nix/store/8bjm87p310sb7r2r0sg4xrynlvg86j8k-hello-2.12.1.tar.gz.drv
```

```
λ docker tag alpine ttl.sh/alpine/alpine:1m
λ docker push ttl.sh/alpine/alpine:1m

skopeo list-tags docker://ttl.sh/alpine/alpine

skopeo copy docker-archive:/tmp/busybox.tar.gz docker://ttl.sh/alpine/alpine
```
