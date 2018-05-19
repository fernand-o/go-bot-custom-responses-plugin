package customresponses

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

const (
	argumentsExample = "Usage: \n !responses set 'Is someone there?' 'Hello' \n !responses unset 'Is someone there?' \n !responses list"
	invalidArguments = "Please inform the params, ex:"
)

var Keys []string
var RedisClient *redis.Client

func connectRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://:@localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	RedisClient = redis.NewClient(opt)
}

func loadMessages() {
	var err error
	Keys, err = RedisClient.Keys("*").Result()
	if err != nil {
		panic(err)
	}
}

func setResponse(args []string) string {
	if (args[0] != "set") || (args[1] == "") || (args[2] == "") {
		return argumentsExample
	}
	match := args[1]
	response := args[2]
	err := RedisClient.Set(match, response, 0).Err()
	if err != nil {
		panic(err)
	}
	loadMessages()
	return confirmationMessageSetResponse(match, response)
}

func getResponse(key string) string {
	response, err := RedisClient.Get(key).Result()
	if err != nil {
		panic(err)
	}
	return response
}

func confirmationMessageSetResponse(match string, response string) string {
	return fmt.Sprintf("Ok! I will send a message with %s when i found any occurences of %s", response, match)
}

func confirmationMessageUnsetResponse(match string) string {
	return fmt.Sprintf("Done, i'll not say anything more related to %s", match)
}

func listResponses(param string) string {
	if param != "list" {
		return argumentsExample
	}
	return "listing responses.."
}

func unsetResponse(param, match string) string {
	if (param != "unset") || (match == "") {
		return argumentsExample
	}
	// (remove from redis)
	return confirmationMessageUnsetResponse(match)
}

func responsesCommand(command *bot.Cmd) (msg string, err error) {
	switch len(command.Args) {
	case 1:
		msg = listResponses(command.Args[0])
	case 2:
		msg = unsetResponse(command.Args[0], command.Args[1])
	case 3:
		msg = setResponse(command.Args)
	default:
		msg = argumentsExample
	}
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
	connectRedis()
	bot.RegisterPassiveCommand(
		"customresponses",
		customresponses)
	bot.RegisterCommand(
		"responses",
		"Defines a custom response to be sent when a given string is found in a message",
		argumentsExample,
		responsesCommand)
	loadMessages()
}
