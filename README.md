# kbot

A small Telegram bot written in Go with [cobra](https://github.com/spf13/cobra) and [telebot v4](https://gopkg.in/telebot.v4).

**Try it:** [t.me/dtkbotedu_bot](https://t.me/dtkbotedu_bot)

## What it does

| You send | Kbot answers |
|---|---|
| `/start` | greeting and the list of commands |
| `/help` | the list of commands |
| `/version` | the bot version, e.g. `kbot v1.0.0` |
| `/ping` | `pong` |
| `hello`, `hi`, `привіт` | a greeting with the version |
| `?` | the list of commands |
| any other text | echoes it back: `You said: …` |
| a photo | its size in pixels |
| a sticker | its emoji |
| a file | its name and size |
| a voice message | its duration |
| a location | its coordinates |

## Install

You need Go 1.22+ and a bot token from [@BotFather](https://t.me/BotFather) (`/newbot`).

```bash
git clone https://github.com/exxmanster/kbot.git
cd kbot
make build              # or: go build -o kbot .
```

## Run

The token is read from the `TELE_TOKEN` environment variable. `read -s` keeps it out of your shell history:

```bash
read -s TELE_TOKEN && export TELE_TOKEN
./kbot start
```

Other commands:

```bash
./kbot version          # print the version
./kbot --help           # list commands
make test               # run unit tests
```

## Build for other platforms

Go cross-compiles, so every binary is built natively on any machine, without QEMU or Rosetta:

| Command | Result |
|---|---|
| `make linux` | `dist/kbot-linux-amd64` |
| `make arm` | `dist/kbot-linux-arm64` |
| `make macos` | `dist/kbot-darwin-arm64` (Apple Silicon) |
| `make windows` | `dist/kbot-windows-amd64.exe` |

Any other pair works too: `make build TARGETOS=linux TARGETARCH=386`.

## Docker image

The [Dockerfile](Dockerfile) uses the Go image from Quay (`quay.io/projectquay/golang`) instead of Docker Hub,
and images are tagged for GitHub Container Registry (`ghcr.io/exxmanster/kbot`). It follows the
[cross-compilation](https://docs.docker.com/build/building/multi-platform/#cross-compilation) pattern:
the build stage runs on the host's own platform and Go compiles for the target one, so the binary
in the image is native for the target.

```bash
make image              # image for the platform of this Docker host, e.g. linux/arm64 on an Apple Silicon Mac
make image-linux        # linux/amd64
make image-arm          # linux/arm64, runs natively on Apple Silicon Macs
make image-windows      # windows/amd64 on a nanoserver base, see the note below
make push               # push the image to ghcr.io (after docker login ghcr.io)
make clean              # remove binaries and the image (docker rmi)
```

Tags look like `ghcr.io/exxmanster/kbot:v1.0.0-linux-arm64`. Run the bot from an image:

```bash
docker run --rm -e TELE_TOKEN ghcr.io/exxmanster/kbot:v1.0.0-linux-arm64
```

All targets use plain `docker build --platform <os>/<arch>`, no buildx. The Windows image can be built
on Linux or macOS, but only a Docker engine with the containerd image store can keep it (the default
since Docker Engine 29; in Docker Desktop: Settings → General → "Use containerd for pulling and storing
images"). Push it with `make push IMAGE_OS=windows IMAGE_ARCH=amd64` and run it on a Windows host in
Windows containers mode.

## Project layout

```
main.go          entry point
cmd/root.go      root cobra command
cmd/kbot.go      `kbot start`: creates the bot and registers handlers
cmd/reply.go     reply logic for text, help text, command menu
cmd/version.go   `kbot version`, set at build time via -ldflags
Makefile         build, cross-builds, images, tests
Dockerfile       multi-platform image: Linux on scratch, Windows on nanoserver
```
