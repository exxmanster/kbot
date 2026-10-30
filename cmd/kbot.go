package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v4"
)

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Run the Telegram bot",
	Long: `Starts the bot with long polling. The token is read from the TELE_TOKEN
environment variable.

  export TELE_TOKEN=<token from @BotFather>
  ./kbot start`,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("TELE_TOKEN")
		if token == "" {
			log.Fatal("TELE_TOKEN environment variable is not set")
		}

		bot, err := telebot.NewBot(telebot.Settings{
			Token:  token,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatalf("Please check TELE_TOKEN: %s", err)
		}

		if err := bot.SetCommands(botCommands); err != nil {
			log.Printf("could not set the command menu: %s", err)
		}
		registerHandlers(bot)

		log.Printf("kbot %s started as @%s", appVersion, bot.Me.Username)
		bot.Start()
	},
}

// logIncoming writes every handled message and whether the reply went out to the log.
func logIncoming(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		var from int64
		if s := c.Sender(); s != nil {
			from = s.ID
		}
		err := next(c)
		status := "replied"
		if err != nil {
			status = "reply failed: " + err.Error()
		}
		log.Printf("%s from %d: %q -> %s", messageKind(c.Message()), from, c.Text(), status)
		return err
	}
}

// registerHandlers attaches a handler for every command and message type the bot understands.
// The logger is added first, because telebot wraps each handler with the middleware set at Handle time.
func registerHandlers(bot *telebot.Bot) {
	bot.Use(logIncoming)

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(fmt.Sprintf("Hi, %s! 👋\n\n%s", c.Sender().FirstName, helpText))
	})
	bot.Handle("/help", func(c telebot.Context) error {
		return c.Send(helpText)
	})
	bot.Handle("/version", func(c telebot.Context) error {
		return c.Send("kbot " + appVersion)
	})
	bot.Handle("/ping", func(c telebot.Context) error {
		return c.Send("pong")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return c.Send(textReply(c.Text(), appVersion))
	})
	bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		p := c.Message().Photo
		return c.Send(fmt.Sprintf("Nice photo! 📷 %d×%d px", p.Width, p.Height))
	})
	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		return c.Send("Cool sticker " + c.Message().Sticker.Emoji)
	})
	bot.Handle(telebot.OnDocument, func(c telebot.Context) error {
		d := c.Message().Document
		return c.Send(fmt.Sprintf("Got a file: %s (%s)", d.FileName, humanSize(d.FileSize)))
	})
	bot.Handle(telebot.OnVoice, func(c telebot.Context) error {
		return c.Send(fmt.Sprintf("Voice message, %d s. I can't listen yet 🙉", c.Message().Voice.Duration))
	})
	bot.Handle(telebot.OnLocation, func(c telebot.Context) error {
		l := c.Message().Location
		return c.Send(fmt.Sprintf("You are at %.5f, %.5f 📍", l.Lat, l.Lng))
	})
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
