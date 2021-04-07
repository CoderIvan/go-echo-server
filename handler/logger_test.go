package handler

import (
	"go-echo-server/datagram"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormat(t *testing.T) {
	Convey("TestFormat 01", t, func() {
		data := datagram.Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549354 * 1e6,
		}

		results := format(data)

		So(results, ShouldNotBeNil)
		So(results[0], ShouldEqual, time.Unix(0, data.Time).Format("20060102 15:04:05.000 +0800"))
		So(results[0], ShouldEqual, "20210401 21:25:49.354 +0800")
		So(results[1], ShouldEqual, ">>")
		So(results[2], ShouldEqual, data.TagName)
		So(results[3], ShouldEqual, ">>")
		So(results[4], ShouldEqual, data.Addr)
		So(results[5], ShouldEqual, ">>")
		So(results[6], ShouldEqual, data.Content)
	})

	Convey("TestFormat 02", t, func() {
		data := datagram.Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "newProject",
			Time:        1617283549354 * 1e6,
		}

		results := format(data)

		So(results, ShouldNotBeNil)
		So(results[0], ShouldEqual, time.Unix(0, data.Time).Format("20060102 15:04:05.000 +0800"))
		So(results[0], ShouldEqual, "20210401 21:25:49.354 +0800")
		So(results[1], ShouldEqual, ">>")
		So(results[2], ShouldEqual, data.TagName)
		So(results[3], ShouldEqual, ">>")
		So(results[4], ShouldEqual, data.Addr)
		So(results[5], ShouldEqual, ">>")
		So(results[6], ShouldEqual, data.ProjectName)
		So(results[7], ShouldEqual, ">>")
		So(results[8], ShouldEqual, data.Content)
	})
}
