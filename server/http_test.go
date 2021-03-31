package server

import (
	"bytes"
	"go-echo-server/datagram"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetupRouter(t *testing.T) {
	Convey("TestSetupRouter rule 1", t, func() {
		ch := make(chan datagram.Datagram, 1)

		router := setupRouter(func(data datagram.Datagram) {
			ch <- data
		})

		w := httptest.NewRecorder()
		now := time.Now().UnixNano()
		content := "Some bytes:0.5239698258003584"
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(content)))
		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, 200)
		result := <-ch
		So(result, ShouldNotBeNil)
		// So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, "")
		So(result.TagName, ShouldEqual, "http-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+1e7)
	})

	Convey("TestSetupRouter rule 2", t, func() {
		ch := make(chan datagram.Datagram, 1)

		router := setupRouter(func(data datagram.Datagram) {
			ch <- data
		})

		w := httptest.NewRecorder()
		now := time.Now().UnixNano()
		content := `{"sn":"0000000006","iccid":"90000000000000000006","imei":"900000000000006","random":"Some bytes:0.28632299830570695"}`
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(content)))
		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, 200)
		result := <-ch
		So(result, ShouldNotBeNil)
		// So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, "")
		So(result.TagName, ShouldEqual, "http-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+1e7)
	})

	Convey("TestSetupRouter rule 3", t, func() {
		ch := make(chan datagram.Datagram, 1)

		router := setupRouter(func(data datagram.Datagram) {
			ch <- data
		})

		w := httptest.NewRecorder()
		now := time.Now().UnixNano()
		projectName := "project_99"
		content := `{"sn":"0000000006","iccid":"90000000000000000006","imei":"900000000000006","random":"Some bytes:0.2391908655820043"}`
		req, _ := http.NewRequest("POST", "/"+projectName, bytes.NewBuffer([]byte(content)))
		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, 200)
		result := <-ch
		So(result, ShouldNotBeNil)
		// So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, projectName)
		So(result.TagName, ShouldEqual, "http-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+1e7)
	})
}

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

		now := time.Now().UnixNano()
		content := "Hello World"
		go func() {
			http.Post("http://127.0.0.1:81", "application/json", bytes.NewBuffer([]byte(content)))
		}()

		result := <-ch
		So(result, ShouldNotBeNil)
		// So(result.Addr, ShouldEqual, Addr)
		So(result.Content, ShouldEqual, content)
		So(result.ProjectName, ShouldEqual, "")
		So(result.TagName, ShouldEqual, "http-server")
		So(result.Time, ShouldBeBetweenOrEqual, now, now+1e9)
	})
}
