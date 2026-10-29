# kbot

A small Telegram bot written in Go with [cobra](https://github.com/spf13/cobra) and [telebot v4](https://gopkg.in/telebot.v4).

**Try it:** [t.me/BOT_USERNAME](https://t.me/BOT_USERNAME)

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

## Project layout

```
main.go          entry point
cmd/root.go      root cobra command
cmd/kbot.go      `kbot start`: creates the bot and registers handlers
cmd/reply.go     reply logic for text, help text, command menu
cmd/version.go   `kbot version`, set at build time via -ldflags
Makefile         format, lint, test, build, run
```
