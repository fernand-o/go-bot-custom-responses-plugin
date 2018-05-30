package customresponses

import (
	"strings"
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

		Convey("match", func() {
			Convey("show, clear", func() {
				sendCommandAndAssertMessage([]string{"match", "showall"}, userMessageNoResposesDefined())

				ActiveCmd.Args = []string{"match", "set", "Life meaning", "42"}
				_, _ = responsesCommand(ActiveCmd)
				ActiveCmd.Args = []string{"match", "set", "I don't know Rick", "Just shoot them Morty"}
				_, _ = responsesCommand(ActiveCmd)

				list := "List of defined responses:\n```\nI don't know Rick -> Just shoot them Morty\nLife meaning -> 42\n```"
				sendCommandAndAssertMessage([]string{"match", "showall"}, list)

				sendCommandAndAssertMessage([]string{"match", "clear"}, userMessageResponsesDeleted())
				sendCommandAndAssertMessage([]string{"match", "showall"}, userMessageNoResposesDefined())
			})

			Convey("set, unset", func() {
				match := "Is someone there?"
				response := "Hello"
				ActiveCmd.Args = []string{"match", "set", match, response}
				msg, err := responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageSetResponse(match, response))

				passiveCmd.Raw = "Hey! Is someone there?"
				msg, err = customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, response)

				ActiveCmd.Args = []string{"match", "unset", "0"}
				msg, err = responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageUnsetResponse(match))

				passiveCmd.Raw = "Hey! Is someone there?"
				msg, err = customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "")
			})
		})

		Convey("list", func() {
			Convey("showall, show, delete", func() {
				sendCommandAndAssertMessage([]string{"list", "showall"}, userMessageNoListsDefined())

				funfacts := []string{"Bananas are curved because they grow towards the sun.", "If you lift a kagaroo's tail off the ground it can't hop."}
				sadfacts := []string{"Heart attacks are more likely to happen on a Monday."}

				ActiveCmd.Args = []string{"list", "add", "#funfacts", funfacts[0]}
				_, _ = responsesCommand(ActiveCmd)
				ActiveCmd.Args = []string{"list", "add", "#funfacts", funfacts[1]}
				_, _ = responsesCommand(ActiveCmd)
				ActiveCmd.Args = []string{"list", "add", "#sadfacts", sadfacts[0]}
				_, _ = responsesCommand(ActiveCmd)

				Convey("showall", func() {
					list := strings.Join([]string{
						"Defined lists:",
						"```",
						"#sadfacts",
						" - " + sadfacts[0],
						"",
						"#funfacts",
						" - " + funfacts[1],
						" - " + funfacts[0],
						"",
						"```"}, "\n")
					sendCommandAndAssertMessage([]string{"list", "showall"}, list)
				})

				Convey("show", func() {
					list := strings.Join([]string{
						"```",
						"#funfacts",
						" - " + funfacts[1],
						" - " + funfacts[0],
						"```"}, "\n")
					sendCommandAndAssertMessage([]string{"list", "show", "#funfacts"}, list)
				})

				Convey("delete", func() {
					sendCommandAndAssertMessage([]string{"list", "delete", "#funfacts"}, userMessageListDeleted("#funfacts"))
				})
			})

			Convey("add, remove", func() {
				listname := "#randomfacts"
				message := "You cannot snore and dream at the same time."
				ActiveCmd.Args = []string{"list", "add", listname, message}
				msg, err := responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageListMessageAdded(listname, message))

				ActiveCmd.Args = []string{"list", "remove", listname, message}
				msg, err = responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageListMessageRemoved(listname, message))

				ActiveCmd.Args = []string{"list", "add", "invalidname", message}
				msg, err = responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, userMessageListInvalidName())
			})
		})

		Convey("clearall", func() {
			ActiveCmd.Args = []string{"list", "add", "#dummy", "message"}
			_, _ = responsesCommand(ActiveCmd)
			ActiveCmd.Args = []string{"match", "set", "Life meaning", "42"}
			_, _ = responsesCommand(ActiveCmd)

			ActiveCmd.Args = []string{"list", "showall"}
			msg, _ := responsesCommand(ActiveCmd)
			So(msg, ShouldNotEqual, userMessageNoListsDefined())

			ActiveCmd.Args = []string{"match", "showall"}
			msg, _ = responsesCommand(ActiveCmd)
			So(msg, ShouldNotEqual, userMessageNoResposesDefined())

			ActiveCmd.Args = []string{"clearall"}
			msg, _ = responsesCommand(ActiveCmd)
			So(msg, ShouldEqual, userMessageDBErased())
		})
	})
}
