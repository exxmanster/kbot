package cmd

import (
	"strings"
	"testing"

	telebot "gopkg.in/telebot.v4"
)

func TestTextReply(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"greeting en", "hello", "Hello, I'm Kbot v1.0.0!"},
		{"greeting case and spaces", "  HeLLo  ", "Hello, I'm Kbot v1.0.0!"},
		{"greeting uk", "Привіт", "Привіт! Я Kbot v1.0.0."},
		{"ping", "ping", "pong"},
		{"question mark asks for help", "?", helpText},
		{"anything else is echoed", "deploy on friday", "You said: deploy on friday"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := textReply(tt.in, "v1.0.0"); got != tt.want {
				t.Errorf("textReply(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHelpTextListsCommands(t *testing.T) {
	for _, c := range botCommands {
		if !strings.Contains(helpText, "/"+c.Text) {
			t.Errorf("helpText does not mention /%s", c.Text)
		}
	}
}

func TestMessageKind(t *testing.T) {
	tests := []struct {
		in   *telebot.Message
		want string
	}{
		{nil, "update"},
		{&telebot.Message{Text: "/start"}, "command"},
		{&telebot.Message{Text: "привіт"}, "text"},
		{&telebot.Message{Photo: &telebot.Photo{}, Caption: "cat"}, "photo"},
		{&telebot.Message{Sticker: &telebot.Sticker{}}, "sticker"},
		{&telebot.Message{Document: &telebot.Document{}}, "file"},
		{&telebot.Message{Voice: &telebot.Voice{}}, "voice"},
		{&telebot.Message{Location: &telebot.Location{}}, "location"},
		{&telebot.Message{}, "other"},
	}
	for _, tt := range tests {
		if got := messageKind(tt.in); got != tt.want {
			t.Errorf("messageKind(%+v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestHumanSize(t *testing.T) {
	tests := map[int64]string{
		512:             "512 B",
		2048:            "2.0 KB",
		5 * 1024 * 1024: "5.0 MB",
	}
	for in, want := range tests {
		if got := humanSize(in); got != want {
			t.Errorf("humanSize(%d) = %q, want %q", in, got, want)
		}
	}
}
