package customresponses

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

const (
	argumentsExample = "Usage: \n !responses set \"Is someone there?\" \"Hello\" \n !responses unset \"Is someone there?\" \n !responses list"
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

func loadKeys() {
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
	return userMessageSetResponse(match, response)
}

func getResponse(key string) string {
	response, _ := RedisClient.Get(key).Result()
	return response
}

func userMessageSetResponse(match string, response string) string {
	return fmt.Sprintf("Ok! I will send a message with %s when i found any occurences of %s", response, match)
}

func userMessageUnsetResponse(match string) string {
	return fmt.Sprintf("Done, i'll not say anything more related to %s", match)
}

func userMessageNoResposesDefined() string {
	return fmt.Sprintf("There's no responses defined yet. \n %s", argumentsExample)
}

func userMessageResponsesDeleted() string {
	return "All responses were deleted."
}

func listOrClearResponses(param string) (msg string) {
	switch param {
	case "list":
		msg = listResponses()
	case "clear":
		msg = clearResponses()
	default:
		msg = argumentsExample
	}
	return
}

func clearResponses() string {
	RedisClient.FlushDB()
	return userMessageResponsesDeleted()
}

func listResponses() string {
	if len(Keys) == 0 {
		return userMessageNoResposesDefined()
	}

	var list, line []string
	for _, k := range Keys {
		line = []string{k, getResponse(k)}
		list = append(list, strings.Join(line, " -> "))
	}
	sort.Sort(sort.StringSlice(list))
	list = append([]string{"List of defined responses:"}, list...)
	return strings.Join(list, "\n")
}

func unsetResponse(param, match string) string {
	if (param != "unset") || (match == "") {
		return argumentsExample
	}
	RedisClient.Del(match)
	return userMessageUnsetResponse(match)
}

func responsesCommand(command *bot.Cmd) (msg string, err error) {
	switch len(command.Args) {
	case 1:
		loadKeys()
		msg = listOrClearResponses(command.Args[0])
	case 2:
		msg = unsetResponse(command.Args[0], command.Args[1])
		loadKeys()
	case 3:
		msg = setResponse(command.Args)
		loadKeys()
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
}
