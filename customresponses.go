package customresponses

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

const (
	argumentsExample = "!setReponse 'Found a banana' 'banana*'"
	invalidArguments = "Please inform the params, ex:"
)

var Keys []string

func newClient() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("REDIS_URL env var not defined")
	}
	return redis.NewClient(opt)
}

func loadMessages() {
	var err error
	client := newClient()
	Keys, err = client.Keys("*").Result()
	if err != nil {
		panic(err)
	}
}

func setResponse(pattern string, response string) {
	client := newClient()
	err := client.Set(pattern, response, 0).Err()
	if err != nil {
		panic(err)
	}
	loadMessages()
}

func getResponse(pattern string) string {
	client := newClient()
	response, err := client.Get(pattern).Result()
	if err != nil {
		panic(err)
	}
	return response
}

func responseMessage(response, pattern string) string {
	return fmt.Sprintf("Ok! I will send a message with %s when i found any matches with %s", response, pattern)
}

func setReponseCommand(command *bot.Cmd) (msg string, err error) {
	if len(command.Args) != 2 {
		msg = argumentsExample
		return
	}
	response := command.Args[0]
	pattern := command.Args[1]
	setResponse(pattern, response)
	msg = responseMessage(pattern, response)
	return
}

func customresponses(command *bot.PassiveCmd) (msg string, err error) {
	var match bool
	for _, k := range Keys {
		match, err = regexp.MatchString(k, command.Raw)
		if match {
			msg = getResponse(k)
			break
		}
	}
	return
}

func init() {
	bot.RegisterPassiveCommand(
		"customresponses",
		customresponses)
	bot.RegisterCommand(
		"setResponse",
		"Defines a custom response for the given pattern",
		argumentsExample,
		setReponseCommand)
	loadMessages()
}
