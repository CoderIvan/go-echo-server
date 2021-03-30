package server

import (
	"fmt"
	"go-echo-server/datagram"
	"net"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProcess(t *testing.T) {
	Convey("TestProcess", t, func() {
		strings := [4]string{
			"Some bytes:0.7913967715224008",
			`{"sn":"0000000002","iccid":"90000000000000000002","imei":"900000000000002","random":"Some bytes:0.7365780627042167"}`,
			`$project_52#{"sn":"0000000005","iccid":"90000000000000000005","imei":"900000000000005","random":"Some bytes:0.5045852142476737"}`,
			`imei=866262040014556,iccid=89860411101871328036,project=LedIndicator4GAll,version=4.1.0,gsm=18,led={"C":10000,"D":{"R":-1,"T":71,"F":28,"L":-1}},led2={"C":10000,"D":{"R":-1,"T":176,"F":104,"L":-1}},reboot=getFailA,request=7168`,
		}

		for _, s := range strings {
			data := process([]byte(s), "127.0.0.1")
			So(data, ShouldNotBeNil)
		}
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

		go func() {
			udpConn, err := net.Dial("udp", ":80")
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			_, err = udpConn.Write([]byte("Hello World"))
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}()

		data := <-ch
		So(data, ShouldNotBeNil)
	})
}
