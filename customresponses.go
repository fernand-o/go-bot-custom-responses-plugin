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

func setResponse(pattern string, response string) {
	err := RedisClient.Set(pattern, response, 0).Err()
	if err != nil {
		panic(err)
	}
	loadMessages()
}

func getResponse(pattern string) string {
	response, err := RedisClient.Get(pattern).Result()
	if err != nil {
		panic(err)
	}
	return response
}

func responseMessage(pattern, response string) string {
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
	connectRedis()
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
