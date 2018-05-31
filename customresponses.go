package customresponses

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

type Match struct {
	key      string
	match    string
	response string
	list     string
}

var Matches []Match
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
	Matches = []Match{}
	var err error
	matches, err := RedisClient.Keys(matchesKeyFmt("*")).Result()
	if err != nil {
		panic(err)
	}

	var values map[string]string
	var match Match
	for _, key := range matches {
		values, _ = RedisClient.HGetAll(key).Result()
		match.match = values["match"]
		match.response = values["response"]
		match.list = values["list"]
		match.key = key
		Matches = append(Matches, match)
	}
}

func showOrClearResponses(param string) (msg string) {
	switch param {
	case "showall":
		msg = showResponses()
	case "clear":
		msg = clearResponses()
	default:
		msg = argumentsMatchExample
	}
	return
}

func clearAll() string {
	RedisClient.FlushDB()
	return userMessageDBErased()
}

func clearResponses() string {
	i := 0
	for i <= recordCount() {
		RedisClient.Del(matchesKeyFmt(strconv.Itoa(i)))
		i++
	}
	return userMessageResponsesDeleted()
}

func showResponses() string {
	if len(Matches) == 0 {
		return userMessageNoResposesDefined()
	}

	var results []string
	var line string
	for _, k := range Matches {
		line = fmt.Sprintf("[%s] [%s] [%s] [%s]", k.key, k.match, k.response, k.list)
		results = append(results, line)
	}
	sort.Sort(sort.StringSlice(results))
	return fmt.Sprintf("List of defined responses:\n```\n[key] [match] [response] [list]\n%s\n```", strings.Join(results, "\n"))
}

func recordCount() int {
	count, _ := RedisClient.DBSize().Result()
	return int(count)
}

func setResponse(args []string) string {
	if args[0] != "set" {
		return argumentsMatchExample
	}

	match := args[1]
	response := args[2]
	list := "_"
	if len(args) == 4 {
		list = args[3]
	}

	params := map[string]interface{}{
		"match":    match,
		"response": response,
		"list":     list}

	count := recordCount()
	key := matchesKeyFmt(strconv.Itoa(count))
	err := RedisClient.HMSet(key, params).Err()
	if err != nil {
		panic(err)
	}
	return userMessageSetResponse(match, response)
}

func unsetResponse(param, id string) string {
	if param != "unset" {
		return argumentsMatchExample
	}
	key := matchesKeyFmt(id)

	_, err := RedisClient.Del(key).Result()
	if err != nil {
		panic(err)
	}
	for _, m := range Matches {
		if m.key == key {
			return userMessageUnsetResponse(m.match)
		}
	}
	return ""
}

func matchCommand(args []string) (msg string) {
	switch len(args) {
	case 1:
		loadMatches()
		msg = showOrClearResponses(args[0])
	case 2:
		msg = unsetResponse(args[0], args[1])
		loadMatches()
	case 3, 4:
		msg = setResponse(args)
		loadMatches()
	default:
		msg = argumentsMatchExample
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
		results = append(results, fmt.Sprintf("  [%s]", m))
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
	listname := args[1]
	message := args[2]
	if !strings.HasPrefix(listname, "#") {
		return userMessageListInvalidName()
	}

	switch args[0] {
	case "add":
		return addListMessage(listname, message)
	case "remove":
		return removeListMessage(listname, message)
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
		msg = argumentsListExample
	}
	return
}

func responsesCommand(command *bot.Cmd) (msg string, err error) {
	paramCount := len(command.Args)
	if paramCount == 0 {
		msg = argumentsGeneralExample
		return
	}

	operation := command.Args[0]

	if (paramCount < 2) && (operation != "clearall") {
		msg = argumentsGeneralExample
		return
	}

	args := append([]string{}, command.Args[1:]...)

	switch operation {
	case "match":
		msg = matchCommand(args)
	case "list":
		msg = listCommand(args)
	case "clearall":
		msg = clearAll()
	default:
		msg = argumentsGeneralExample
	}
	return
}

func getListMessage(listname string) string {
	msg, _ := RedisClient.SRandMember(listname).Result()
	return msg
}

func getFormattedMessage(response, listname string) string {
	if listname == "" {
		return response
	}

	message := getListMessage(listname)
	if strings.Contains(message, "%s") {
		return fmt.Sprintf(message, response)
	}

	return message + response
}

func customresponses(command *bot.PassiveCmd) (msg string, err error) {
	var match bool
	for _, m := range Matches {
		match, err = regexp.MatchString(m.match, command.Raw)
		if match {
			msg = getFormattedMessage(m.response, m.list)
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
		argumentsGeneralExample,
		responsesCommand)
}
