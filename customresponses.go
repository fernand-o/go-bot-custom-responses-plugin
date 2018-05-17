package customresponses

import (
	"github.com/go-chat-bot/bot"
)

func customresponses(command *bot.PassiveCmd) (string, error) {
	return "@fernando", nil
}

func init() {
	bot.RegisterPassiveCommand(
		"customresponses",
		customresponses)
}
