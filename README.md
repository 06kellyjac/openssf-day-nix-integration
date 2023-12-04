https://nixos.org/manual/nix/stable/advanced-topics/post-build-hook

```
nix-store --generate-binary-cache-key cache-key ./nix.key ./nix.pub
```

```
./nix-closure-graph/nix-closure-graph --lg ./result --recursive --derivation > lg.json

nix-store --query --graph $(nix-store --query --deriver ./result) | tred | dot -Tsvg > graphhh.svg
```

```
docker compose ls --filter "name=guac"
guacone collect files guac-data-main/docs

curl 'http://localhost:8080/query' -s -X POST -H 'content-type: application/json' \
  --data '{
    "query": "{ packages(pkgSpec: {}) { type } }"
  }' | jq


sbomnix ./result --type buildtime --spdx
sed -i 's/pkg:nix/pkg:generic/g' sbom.cdx.json

nixgraph ./result --depth=10 --buildtime --colorize 'hello-mini|hello\.c|bash-.*\d\.drv|bash-.*\db\.drv|kaem-\d.*\d\.drv|kaem\..*-\d.*\d\.drv|hex2-.*\d\.drv|hex1-.*\d\.drv|hex0-.*\d\.drv|M1-.*\d\.drv|M2-.*\d\.drv|blood-elf-.*\d\.drv|mescc-tools-\d.*\d\.drv|mescc-tools-extra-\d.*\d\.drv|coreutils-.*\d\.drv|musl-.*\d\.drv|.*coreutils-.*\d\.drv|tinycc-musl.*-.*\d[\-compiler|\-libs]+\.drv|tinycc-mes.*-.*\d[\-compiler|\-libs]+\.drv'
```

```
nix path-info /nix/store/ijq43zaps05w99zgxx2fsx2cp5wdhw3y-trivy-0.46.0 --store https://cache.nixos.org/
nix path-info /nix/store/ijq43zaps05w99zgxx2fsx2cp5wdhw3y-trivy-0.46.0 --json | jq
nix path-info /nix/store/maax53ljrr75q29swzsi9as64di9iykz-trivy-0.46.0.drv --json | jq

curl cache.nixos.org/ijq43zaps05w99zgxx2fsx2cp5wdhw3y.narinfo
```

`filter: blur(4px);`

```
λ docker build . -t hello

λ docker run -it hello
Hello, world!
```

```
λ nix log --help


λ nix log /nix/store/vhdj1wl07gljc6q4pc35bygfp6p131dm-hello-mini-0.0.1.drv

λ stat /nix/store/vhdj1wl07gljc6q4pc35bygfp6p131dm-hello-mini-0.0.1.drv
```
