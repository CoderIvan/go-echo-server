package handler

import (
	"go-echo-server/datagram"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSLSProcess(t *testing.T) {
	Convey("TestProcess 01", t, func() {
		data := datagram.Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "",
			Time:        1617283549354 * 1e6,
		}
		logGroup := process(data)

		So(logGroup, ShouldNotBeNil)
		So(len(logGroup.Logs), ShouldEqual, 1)
		logs := logGroup.Logs[0]
		So(*logs.Time, ShouldEqual, data.Time/1e9)
		So(len(logs.Contents), ShouldEqual, 3)

		So(*logs.Contents[0].Key, ShouldEqual, "tagName")
		So(*logs.Contents[0].Value, ShouldEqual, data.TagName)
		So(*logs.Contents[1].Key, ShouldEqual, "content")
		So(*logs.Contents[1].Value, ShouldEqual, data.Content)
		So(*logs.Contents[2].Key, ShouldEqual, "addr")
		So(*logs.Contents[2].Value, ShouldEqual, data.Addr)
	})

	Convey("TestProcess 02", t, func() {
		data := datagram.Datagram{
			TagName:     "udp-server",
			Content:     "Hello World",
			Addr:        "127.0.0.1",
			ProjectName: "newProject",
			Time:        1617283549354 * 1e6,
		}
		logGroup := process(data)

		So(logGroup, ShouldNotBeNil)
		So(len(logGroup.Logs), ShouldEqual, 1)
		logs := logGroup.Logs[0]
		So(*logs.Time, ShouldEqual, data.Time/1e9)
		So(len(logs.Contents), ShouldEqual, 4)

		So(*logs.Contents[0].Key, ShouldEqual, "tagName")
		So(*logs.Contents[0].Value, ShouldEqual, data.TagName)
		So(*logs.Contents[1].Key, ShouldEqual, "content")
		So(*logs.Contents[1].Value, ShouldEqual, data.Content)
		So(*logs.Contents[2].Key, ShouldEqual, "addr")
		So(*logs.Contents[2].Value, ShouldEqual, data.Addr)
		So(*logs.Contents[3].Key, ShouldEqual, "projectName")
		So(*logs.Contents[3].Value, ShouldEqual, data.ProjectName)
	})
}
