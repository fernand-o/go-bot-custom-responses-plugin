package customresponses

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

var ActiveCmd = &bot.Cmd{}

func defineAndTestResponses(pattern, response string) {
	ActiveCmd.Args = []string{response, pattern}
	s, err := setReponseCommand(ActiveCmd)
	So(err, ShouldBeNil)
	So(s, ShouldEqual, responseMessage(pattern, response))
}

func TestCustomResponses(t *testing.T) {
	Convey("Given a text", t, func() {

		passiveCmd := &bot.PassiveCmd{}
		Convey("When the text doesn't match a defined pattern", func() {
			passiveCmd.Raw = "lorem ipsum"
			s, err := customresponses(passiveCmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text matches a defined pattern", func() {
			defineAndTestResponses("fruits", "@apple")
			passiveCmd.Raw = "The fruits are getting stale"
			s, err := customresponses(passiveCmd)
			So(err, ShouldBeNil)
			So(s, ShouldEqual, "@apple")
		})

		Convey("When the text matches a defined pattern twice", func() {
			defineAndTestResponses("iron-man", "@HomemDeFerro")
			passiveCmd.Raw = "iron-mand: The avengers are stronger with iron-man"
			s, err := customresponses(passiveCmd)
			So(err, ShouldBeNil)
			So(s, ShouldEqual, "@HomemDeFerro")
		})
	})
}
