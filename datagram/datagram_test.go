package datagram

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEqual(t *testing.T) {
	Convey("TestEqual return true", t, func() {
		data01 := Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549354 * 1e6,
		}

		data02 := Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549354 * 1e6,
		}

		So(data01.Equal(data02), ShouldBeTrue)
	})

	Convey("TestEqual return false", t, func() {
		data01 := Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549354 * 1e6,
		}

		data02 := Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549353 * 1e6,
		}

		So(data01.Equal(data02), ShouldBeFalse)
	})
}
