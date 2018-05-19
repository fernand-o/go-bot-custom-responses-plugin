package customresponses

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

var ActiveCmd = &bot.Cmd{}

func defineAndTestResponses(match, response string) {
	ActiveCmd.Args = []string{match, response}
	s, err := responsesCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(s, ShouldEqual, confirmationMessageSetResponse(match, response))
}

func assetInvalidArgument(args []string) {
	ActiveCmd.Args = args
	msg, err := responsesCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(msg, ShouldEqual, argumentsExample)
}

func TestCustomResponses(t *testing.T) {
	Convey("Given a text", t, func() {
		RedisClient.FlushDB()
		// passiveCmd := &bot.PassiveCmd{}

		Convey("!reponses", func() {
			Convey("Wrong parameters", func() {
				assetInvalidArgument([]string{""})
				assetInvalidArgument([]string{"wtf"})
				assetInvalidArgument([]string{"set"})
				assetInvalidArgument([]string{"set", "message"})
				assetInvalidArgument([]string{"set", ""})
				assetInvalidArgument([]string{"unset"})
				assetInvalidArgument([]string{"unset", ""})
				assetInvalidArgument([]string{"list", "wtf"})
			})

			Convey("set", func() {
				match := "Is someone there?"
				response := "Hello"
				ActiveCmd.Args = []string{"set", match, response}
				msg, err := responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, confirmationMessageSetResponse(match, response))
				//test responses
			})

			Convey("unset", func() {
				match := "Is someone there?"
				ActiveCmd.Args = []string{"unset", match}
				msg, err := responsesCommand(ActiveCmd)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, confirmationMessageUnsetResponse(match))
			})
		})

		// Convey("When the text doesn't match a defined pattern", func() {
		// 	passiveCmd.Raw = "lorem ipsum"
		// 	s, err := customresponses(passiveCmd)

		// 	So(err, ShouldBeNil)
		// 	So(s, ShouldEqual, "")
		// })

		// Convey("When the text matches a defined pattern", func() {
		// 	defineAndTestResponses("fruits", "@apple")
		// 	passiveCmd.Raw = "The fruits are getting stale"
		// 	s, err := customresponses(passiveCmd)
		// 	So(err, ShouldBeNil)
		// 	So(s, ShouldEqual, "@apple")
		// })

		// Convey("When the text matches a defined pattern twice", func() {
		// 	defineAndTestResponses("iron-man", "@HomemDeFerro")
		// 	passiveCmd.Raw = "iron-mand: The avengers are stronger with iron-man"
		// 	s, err := customresponses(passiveCmd)
		// 	So(err, ShouldBeNil)
		// 	So(s, ShouldEqual, "@HomemDeFerro")
		// })

	})
}
