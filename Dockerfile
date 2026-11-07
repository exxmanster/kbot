# Multi-platform image for kbot, built with plain `docker build --platform <os>/<arch>`.
#
# The builder stage always runs on the build host's own platform ($BUILDPLATFORM)
# and Go cross-compiles for $TARGETOS/$TARGETARCH. The binary in the final image
# is native for the target, so nothing runs under QEMU or Rosetta.
# See https://docs.docker.com/build/building/multi-platform/#cross-compilation

FROM --platform=$BUILDPLATFORM quay.io/projectquay/golang:1.27 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build TARGETOS=$TARGETOS TARGETARCH=$TARGETARCH VERSION=$VERSION OUT=/out/kbot

# Linux image (amd64, arm64, ...): just the binary and CA certificates for the Telegram API.
FROM scratch AS release-linux
COPY --from=builder /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/kbot /kbot
ENTRYPOINT ["/kbot"]
CMD ["start"]

# Windows image: a Windows container needs a Windows base layer, scratch is not enough.
FROM mcr.microsoft.com/windows/nanoserver:ltsc2022 AS release-windows
COPY --from=builder /out/kbot /kbot.exe
ENTRYPOINT ["C:\\kbot.exe"]
CMD ["start"]

# Pick the final stage by the target OS: release-linux or release-windows.
FROM release-$TARGETOS
