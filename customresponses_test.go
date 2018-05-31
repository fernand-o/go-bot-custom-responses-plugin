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
	So(msg, ShouldEqual, argumentsGeneralExample)
}

func sendCommandAndAssertMessage(args []string, expectedMessage string) {
	ActiveCmd.Args = args
	msg, err := responsesCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(msg, ShouldEqual, expectedMessage)
}

func sendCommand(args []string) {
	ActiveCmd.Args = args
	_, _ = responsesCommand(ActiveCmd)
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
				sendCommand([]string{"match", "set", "Life meaning", "42"})
				sendCommand([]string{"match", "set", "I don't know Rick", "Just shoot them Morty"})
				list := "List of defined responses:\n```\nI don't know Rick -> Just shoot them Morty\nLife meaning -> 42\n```"
				sendCommandAndAssertMessage([]string{"match", "showall"}, list)
				sendCommandAndAssertMessage([]string{"match", "clear"}, userMessageResponsesDeleted())
				sendCommandAndAssertMessage([]string{"match", "showall"}, userMessageNoResposesDefined())
			})

			Convey("set, unset", func() {
				match := "Is someone there?"
				response := "Hello"
				sendCommandAndAssertMessage([]string{"match", "set", match, response}, userMessageSetResponse(match, response))

				passiveCmd.Raw = "Hey! Is someone there?"
				msg, err := customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, response)

				sendCommandAndAssertMessage([]string{"match", "unset", "0"}, userMessageUnsetResponse(match))

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

				sendCommand([]string{"list", "add", "#funfacts", funfacts[0]})
				sendCommand([]string{"list", "add", "#funfacts", funfacts[1]})
				sendCommand([]string{"list", "add", "#sadfacts", sadfacts[0]})

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
				sendCommandAndAssertMessage([]string{"list", "add", listname, message}, userMessageListMessageAdded(listname, message))
				sendCommandAndAssertMessage([]string{"list", "remove", listname, message}, userMessageListMessageRemoved(listname, message))
				sendCommandAndAssertMessage([]string{"list", "add", "invalidname", message}, userMessageListInvalidName())
			})
		})

		Convey("clearall", func() {
			sendCommand([]string{"list", "add", "#dummy", "message"})
			sendCommand([]string{"match", "set", "Life meaning", "42"})

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

		Convey("passive command", func() {
			Convey("with several list messages", func() {
				var possibleResults []string
				msgs := []string{"wubba lubba dub dub", "aw geez Rick", "i don't know, maybe, "}
				response := "lost in space"
				listname := "#dummy"

				for _, m := range msgs {
					sendCommand([]string{"list", "add", listname, m})
					possibleResults = append(possibleResults, m+response)
				}
				sendCommand([]string{"match", "set", "where is my portal gun?", response, listname})

				passiveCmd.Raw = "where is my portal gun?"
				msg, err := customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldBeIn, possibleResults)
			})

			Convey("with formatted message", func() {
				sendCommand([]string{"list", "add", "#fun", "Did you notice that %s is drunk?"})
				sendCommand([]string{"match", "set", "fernand-o is drinking", "@fernando", "#fun"})

				passiveCmd.Raw = "fernand-o is drinking"
				msg, err := customresponses(passiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "Did you notice that @fernando is drunk?")
			})
		})
	})
}
