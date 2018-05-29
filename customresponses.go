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
	argumentsExample     = "Usage: \n```\n!responses set \"Is someone there?\" \"Hello\" \n !responses unset \"Is someone there?\" \n !responses list\n```"
	argumentsListExample = "Usage: \n```\n !responses list add mylist \"Some random message\" \n !responses list delete mylist \"Some random message\" \n !responses list clear mylist\n```"
	invalidArguments     = "Please inform the params, ex:"
	matchesKey           = "matches"
)

var Matches []string
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

func loadMatches() {
	var err error
	Matches, err = RedisClient.HKeys(matchesKey).Result()
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
	err := RedisClient.HSet(matchesKey, match, response).Err()
	if err != nil {
		panic(err)
	}
	return userMessageSetResponse(match, response)
}

func getResponse(key string) string {
	response, _ := RedisClient.HGet(matchesKey, key).Result()
	return response
}

func userMessageSetResponse(match, response string) string {
	return fmt.Sprintf("Ok! I will send a message with %s when i found any occurences of %s", response, match)
}

func userMessageUnsetResponse(match string) string {
	return fmt.Sprintf("Done, i'll not say anything more related to %s", match)
}

func userMessageNoResposesDefined() string {
	return fmt.Sprintf("There are no responses defined yet. \n %s", argumentsExample)
}

func userMessageResponsesDeleted() string {
	return "All responses were deleted."
}

func userMessageListMessageAdded(list, message string) string {
	return fmt.Sprintf("The message `%s` was added to the list `%s`.", message, list)
}

func userMessageListMessageRemoved(list, message string) string {
	return fmt.Sprintf("The message `%s` was removed of the list `%s`.", message, list)
}

func userMessageListDeleted(list string) string {
	return fmt.Sprintf("The list %s was deleted.", list)
}

func userMessageNoListsDefined() string {
	return fmt.Sprintf("There are no lists defined yet. \n %s", argumentsListExample)
}

func showOrClearResponses(param string) (msg string) {
	switch param {
	case "show":
		msg = showResponses()
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

func showResponses() string {
	if len(Matches) == 0 {
		return userMessageNoResposesDefined()
	}

	var results, line []string
	for _, k := range Matches {
		line = []string{k, getResponse(k)}
		results = append(results, strings.Join(line, " -> "))
	}
	sort.Sort(sort.StringSlice(results))
	return fmt.Sprintf("List of defined responses:\n```\n%s\n```", strings.Join(results, "\n"))
}

func unsetResponse(param, match string) string {
	if (param != "unset") || (match == "") {
		return argumentsExample
	}
	RedisClient.HDel(matchesKey, match)
	return userMessageUnsetResponse(match)
}

func matchCommand(args []string) (msg string) {
	switch len(args) {
	case 1:
		loadMatches()
		msg = showOrClearResponses(args[0])
	case 2:
		msg = unsetResponse(args[0], args[1])
		loadMatches()
	case 3:
		msg = setResponse(args)
		loadMatches()
	default:
		msg = argumentsExample
	}
	return
}

func showOrClearList(args []string) string {
	switch args[0] {
	case "show":
		return "```\n" + getListMembers(args[1]) + "\n```"
	case "delete":
		return userMessageListDeleted(args[1])
	default:
		return argumentsListExample
	}
}

func getListMembers(listname string) string {
	var results = []string{listname}
	messages, _ := RedisClient.SMembers(listname).Result()
	for _, m := range messages {
		results = append(results, " - "+m)
	}
	return strings.Join(results, "\n")
}

func showAllLists(param string) string {
	if param != "showall" {
		return argumentsListExample
	}

	lists, _ := RedisClient.Keys("#*").Result()
	if len(lists) == 0 {
		return userMessageNoListsDefined()
	}

	var results []string
	for _, k := range lists {
		results = append(results, getListMembers(k))
		results = append(results, "")
	}

	return fmt.Sprintf("Defined lists:\n```\n%s\n```", strings.Join(results, "\n"))
}

func addListMessage(listname, message string) string {
	err := RedisClient.SAdd(listname, message).Err()
	if err != nil {
		panic(err)
	}
	return userMessageListMessageAdded(listname, message)
}

func removeListMessage(listname, message string) string {
	err := RedisClient.SRem(listname, message).Err()
	if err != nil {
		panic(err)
	}
	return userMessageListMessageRemoved(listname, message)
}

func addOrRemoveListMessage(args []string) string {
	switch args[0] {
	case "add":
		return addListMessage(args[1], args[2])
	case "remove":
		return removeListMessage(args[1], args[2])
	default:
		return argumentsListExample
	}
}

func listCommand(args []string) (msg string) {
	switch len(args) {
	case 1:
		msg = showAllLists(args[0])
	case 2:
		msg = showOrClearList(args)
	case 3:
		msg = addOrRemoveListMessage(args)
	default:
		msg = argumentsExample
	}
	return
}

func responsesCommand(command *bot.Cmd) (msg string, err error) {
	if len(command.Args) < 2 {
		msg = argumentsExample
		return
	}

	operation := command.Args[0]
	args := append([]string{}, command.Args[1:]...)

	switch operation {
	case "match":
		msg = matchCommand(args)
	case "list":
		msg = listCommand(args)
	default:
		msg = argumentsExample
	}
	return
}

func customresponses(command *bot.PassiveCmd) (msg string, err error) {
	var match bool
	for _, k := range Matches {
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
