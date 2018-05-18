package customresponses

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomResponses(t *testing.T) {
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}
		Convey("When the text doesn't match a defined pattern", func() {
			cmd.Raw = "fernand-o"
			s, err := customresponses(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "@fernando")
		})

		Convey("When the text matches a defined pattern", func() {
			cmd.Raw = "[Error] Something went wrong fernand-o, take a look"

			s, err := customresponses(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("Set and Get a response", func() {
			cmd.Raw = ""
			setResponse("iron-man", "@HomemDeFerro")
			s := getResponse("iron-man")
			So(s, ShouldEqual, "@HomemDeFerro")
		})

	})
}
