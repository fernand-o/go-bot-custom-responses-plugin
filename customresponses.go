package customresponses

import (
	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

var Keys []string

func newClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
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
}

func getResponse(pattern string) string {
	client := newClient()
	response, err := client.Get(pattern).Result()
	if err != nil {
		panic(err)
	}
	return response
}

func customresponses(command *bot.PassiveCmd) (string, error) {
	for _, k := range Keys {
		if k == command.Raw {
			return getResponse(k), nil
		}
	}
	return "", nil
}

func init() {
	bot.RegisterPassiveCommand(
		"customresponses",
		customresponses)
	loadMessages()
}
