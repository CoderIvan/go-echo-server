package server

import (
	"go-echo-server/datagram"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessLed(t *testing.T) {
	Convey("TestProcessLed rule 4", t, func() {
		origin := `imei=866262040003450,iccid=89860411101871328041,project=LedIndicator4GLINGYI,version=7.0.5,gsm=24,led={"C":10000,"D":{"R":-1,"F":5,"L":-1,"T":10}},reboot=timer,request=5`

		ok, projectName, content := processLed([]byte(origin))

		So(ok, ShouldBeTrue)
		So(projectName, ShouldEqual, "LedIndicator4GLINGYI")
		So(string(content), ShouldEqual, `{"gsm":"24","iccid":"89860411101871328041","imei":"866262040003450","led":{"C":10000,"D":{"F":5,"L":-1,"R":-1,"T":10}},"reboot":"timer","request":"5","version":"7.0.5"}`)
	})

	Convey("TestProcessLed rule 5", t, func() {
		origin := `imei=866262040014556,iccid=89860411101871328036,project=LedIndicator4GAll,version=4.1.0,gsm=18,led={"C":10000,"D":{"R":-1,"T":71,"F":2,"L":-1}},led2={"C":10000,"D":{"R":-1,"T":176,"F":6,"L":-1}},reboot=timer,request=6272`

		ok, projectName, content := processLed([]byte(origin))

		So(ok, ShouldBeTrue)
		So(projectName, ShouldEqual, "LedIndicator4GAll")
		So(string(content), ShouldEqual, `{"gsm":"18","iccid":"89860411101871328036","imei":"866262040014556","led":{"C":10000,"D":{"F":2,"L":-1,"R":-1,"T":71}},"led2":{"C":10000,"D":{"F":6,"L":-1,"R":-1,"T":176}},"reboot":"timer","request":"6272","version":"4.1.0"}`)
	})
}

func TestProcess(t *testing.T) {
	Convey("TestProcess rule 1", t, func() {
		Addr := "127.0.0.1"
		content := "Some bytes:0.7913967715224008"

		now := time.Now().UnixNano()

		So(process([]byte(content), Addr), ShouldEqual, datagram.Datagram{
			Addr:        Addr,
			Content:     content,
			ProjectName: "",
			TagName:     "udp-server",
			Time:        now,
		})
	})

	Convey("TestProcess rule 2", t, func() {
		Addr := "127.0.0.1"
		content := `{"sn":"0000000002","iccid":"90000000000000000002","imei":"900000000000002","random":"Some bytes:0.7365780627042167"}`

		now := time.Now().UnixNano()

		So(process([]byte(content), Addr), ShouldEqual, datagram.Datagram{
			Addr:        Addr,
			Content:     content,
			ProjectName: "",
			TagName:     "udp-server",
			Time:        now,
		})
	})

	Convey("TestProcess rule 3", t, func() {
		Addr := "127.0.0.1"

		projectName := "project_52"
		content := `{"sn":"0000000005","iccid":"90000000000000000005","imei":"900000000000005","random":"Some bytes:0.5045852142476737"}`

		now := time.Now().UnixNano()

		So(process([]byte("$"+projectName+"#"+content), Addr), ShouldEqual, datagram.Datagram{
			Addr:        Addr,
			Content:     content,
			ProjectName: projectName,
			TagName:     "udp-server",
			Time:        now,
		})
	})

	Convey("TestProcess rule 4", t, func() {
		Addr := "127.0.0.1"

		projectName := "LedIndicator4GLINGYI"
		content := `{"gsm":"24","iccid":"89860411101871328041","imei":"866262040003450","led":{"C":10000,"D":{"F":5,"L":-1,"R":-1,"T":10}},"reboot":"timer","request":"5","version":"7.0.5"}`

		now := time.Now().UnixNano()

		result := process([]byte(`imei=866262040003450,iccid=89860411101871328041,project=LedIndicator4GLINGYI,version=7.0.5,gsm=24,led={"C":10000,"D":{"R":-1,"F":5,"L":-1,"T":10}},reboot=timer,request=5`), Addr)

		So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, projectName)
		So(result.TagName, ShouldEqual, "udp-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+int64(2*time.Millisecond))
	})

	Convey("TestProcess rule 5", t, func() {
		Addr := "127.0.0.1"

		projectName := "LedIndicator4GAll"
		content := `{"gsm":"18","iccid":"89860411101871328036","imei":"866262040014556","led":{"C":10000,"D":{"F":2,"L":-1,"R":-1,"T":71}},"led2":{"C":10000,"D":{"F":6,"L":-1,"R":-1,"T":176}},"reboot":"timer","request":"6272","version":"4.1.0"}`

		now := time.Now().UnixNano()

		result := process([]byte(`imei=866262040014556,iccid=89860411101871328036,project=LedIndicator4GAll,version=4.1.0,gsm=18,led={"C":10000,"D":{"R":-1,"T":71,"F":2,"L":-1}},led2={"C":10000,"D":{"R":-1,"T":176,"F":6,"L":-1}},reboot=timer,request=6272`), Addr)

		So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, projectName)
		So(result.TagName, ShouldEqual, "udp-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+int64(time.Millisecond))
	})
}

func TestUDP(t *testing.T) {
	Convey("TestUDP", t, func() {
		ch := make(chan datagram.Datagram)

		server := NewUDPServer(80)
		defer server.Close()

		go func() {
			server.Listen(func(data datagram.Datagram) {
				ch <- data
			})
		}()

		now := time.Now().UnixNano()
		content := "Hello World"
		var Addr string
		go func() {
			udpConn, _ := net.Dial("udp", ":80")
			Addr = udpConn.LocalAddr().String()
			udpConn.Write([]byte(content))
		}()

		result := <-ch
		So(result, ShouldNotBeNil)
		So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, "")
		So(result.TagName, ShouldEqual, "udp-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+int64(2*time.Millisecond))
	})
}
