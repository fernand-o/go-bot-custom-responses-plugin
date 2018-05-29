package customresponses

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

var ActiveCmd = &bot.Cmd{}

func assetInvalidArgument(args []string) {
	ActiveCmd.Args = args
	msg, err := responsesCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(msg, ShouldEqual, argumentsExample)
}

func sendCommandAndAssertMessage(args []string, expectedMessage string) {
	ActiveCmd.Args = args
	msg, err := responsesCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(msg, ShouldEqual, expectedMessage)
}

func TestCustomResponses(t *testing.T) {
	Convey("Given a text", t, func() {
		RedisClient.FlushDB()
		passiveCmd := &bot.PassiveCmd{}

		Convey("!reponses", func() {
			Convey("Wrong parameters", func() {
				assetInvalidArgument([]string{""})
				assetInvalidArgument([]string{"wtf"})
				assetInvalidArgument([]string{"set"})
				assetInvalidArgument([]string{"set", "message"})
				assetInvalidArgument([]string{"set", ""})
				assetInvalidArgument([]string{"unset"})
				assetInvalidArgument([]string{"unset", ""})
				assetInvalidArgument([]string{"show", "wtf"})
			})

			Convey("show, clear", func() {
				sendCommandAndAssertMessage([]string{"show"}, userMessageNoResposesDefined())

				ActiveCmd.Args = []string{"set", "Life meaning", "42"}
				_, _ = responsesCommand(ActiveCmd)
				ActiveCmd.Args = []string{"set", "I don't know Rick", "Just shoot them Morty"}
				_, _ = responsesCommand(ActiveCmd)

				list := "List of defined responses:\n```\nI don't know Rick -> Just shoot them Morty\nLife meaning -> 42\n```"
				sendCommandAndAssertMessage([]string{"show"}, list)

				sendCommandAndAssertMessage([]string{"clear"}, userMessageResponsesDeleted())
				sendCommandAndAssertMessage([]string{"show"}, userMessageNoResposesDefined())
			})

			Convey("set, unset", func() {
				match := "Is someone there?"
				response := "Hello"
				ActiveCmd.Args = []string{"set", match, response}
				msg, err := responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageSetResponse(match, response))

				passiveCmd.Raw = "Hey! Is someone there?"
				msg, err = customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, response)

				ActiveCmd.Args = []string{"unset", match}
				msg, err = responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageUnsetResponse(match))

				passiveCmd.Raw = "Hey! Is someone there?"
				msg, err = customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "")
			})
		})
	})
}
