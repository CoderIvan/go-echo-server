package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-echo-server/datagram"
	"net/http"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTP(t *testing.T) {
	Convey("TestHTTP", t, func() {
		ch := make(chan datagram.Datagram)
		server := NewHTTPServer(81)
		defer server.Close()

		go func() {
			server.Listen(func(data datagram.Datagram) {
				ch <- data
			})
		}()

		go func() {
			values := map[string]string{"id": "8"}
			jsonValue, _ := json.Marshal(values)
			_, err := http.Post("http://127.0.0.1:81", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}()

		data := <-ch
		So(data, ShouldNotBeNil)
	})
}

func TestHTTP2(t *testing.T) {
	Convey("TestHTTP", t, func() {
		ch := make(chan datagram.Datagram)
		server := NewHTTPServer(81)
		defer server.Close()

		go func() {
			server.Listen(func(data datagram.Datagram) {
				ch <- data
			})
		}()

		go func() {
			values := map[string]string{"id": "8"}
			jsonValue, _ := json.Marshal(values)
			_, err := http.Post("http://127.0.0.1:81/test", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}()

		data := <-ch
		So(data, ShouldNotBeNil)
	})
}
