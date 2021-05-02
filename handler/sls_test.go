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

		testDataList := [3][2]string{
			{"tagName", data.TagName},
			{"addr", data.Addr},
			{"content", data.Content},
		}

		for index, testData := range testDataList {
			So(*logs.Contents[index].Key, ShouldEqual, testData[0])
			So(*logs.Contents[index].Value, ShouldEqual, testData[1])
		}
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

		testDataList := [4][2]string{
			{"tagName", data.TagName},
			{"addr", data.Addr},
			{"projectName", data.ProjectName},
			{"content", data.Content},
		}

		for index, testData := range testDataList {
			So(*logs.Contents[index].Key, ShouldEqual, testData[0])
			So(*logs.Contents[index].Value, ShouldEqual, testData[1])
		}
	})
}
