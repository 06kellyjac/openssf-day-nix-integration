FROM nixos/nix:2.19.2 AS builder
COPY . /files
WORKDIR /files
RUN \
    nix build \
    --extra-experimental-features "nix-command flakes" \
    --print-build-logs \
    /files#hello-static
FROM scratch
COPY --from=builder /files/result/bin/hello /
CMD [ "/hello" ]
