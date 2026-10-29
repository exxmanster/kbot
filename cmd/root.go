package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbot",
	Short: "Kbot is a small Telegram bot written in Go",
	Long: `Kbot is a Telegram bot built with cobra and telebot.

It answers commands (/start, /help, /version, /ping) and reacts to text,
photos, stickers, files, voice messages and locations.

Run "kbot start" with TELE_TOKEN set to launch it.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
