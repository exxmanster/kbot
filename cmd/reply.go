package cmd

import (
	"fmt"
	"strings"

	telebot "gopkg.in/telebot.v4"
)

// botCommands is the menu shown in Telegram next to the input field.
var botCommands = []telebot.Command{
	{Text: "start", Description: "Say hi and show what I can do"},
	{Text: "help", Description: "List commands"},
	{Text: "version", Description: "Show bot version"},
	{Text: "ping", Description: "Check that the bot is alive"},
}

const helpText = `I'm Kbot, a small Telegram bot written in Go.

Commands:
/start   – say hi
/help    – this list
/version – bot version
/ping    – pong

Send me text, a photo, a sticker, a file, a voice message or a location – I'll tell you what I got.`

// textReply decides how the bot answers a plain text message.
func textReply(text, version string) string {
	t := strings.ToLower(strings.TrimSpace(text))
	switch t {
	case "hello", "hi":
		return fmt.Sprintf("Hello, I'm Kbot %s!", version)
	case "привіт":
		return fmt.Sprintf("Привіт! Я Kbot %s.", version)
	case "ping":
		return "pong"
	case "?":
		return helpText
	}
	return "You said: " + strings.TrimSpace(text)
}

// humanSize formats a file size in bytes for chat replies.
func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGT"[exp])
}
