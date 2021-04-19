package server

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func getPackage(stringMsg string) []byte {
	msg := []byte(stringMsg)
	tmp := uint16(len(msg))
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	return append(bytesBuffer.Bytes(), msg...)
}

func TestParserFactory(t *testing.T) {
	Convey("正常情况", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		parser(getPackage("Hello World"))

		So(len(results), ShouldEqual, 1)
		So(string(results[0]), ShouldEqual, "Hello World")
	})

	Convey("正常情况，多个包", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		parser(getPackage("Hello"))
		parser(getPackage("World"))

		So(len(results), ShouldEqual, 2)
		So(string(results[0]), ShouldEqual, "Hello")
		So(string(results[1]), ShouldEqual, "World")
	})

	Convey("两个包，粘包", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		pkg01 := getPackage("The news rocked the world.")
		pkg02 := getPackage("The tragedy resounded around the world.")

		parser(append(pkg01, pkg02...))

		So(len(results), ShouldEqual, 2)
		So(string(results[0]), ShouldEqual, "The news rocked the world.")
		So(string(results[1]), ShouldEqual, "The tragedy resounded around the world.")
	})

	Convey("一个包，分包", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		pkg := getPackage("The news rocked the world.")
		rNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(pkg))
		parser(pkg[:rNum])
		parser(pkg[rNum:])

		So(len(results), ShouldEqual, 1)
		So(string(results[0]), ShouldEqual, "The news rocked the world.")
	})

	Convey("两个包，前面一个包被分包，后半部分和后面的粘包", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		pkg01 := getPackage("The news rocked the world.")
		pkg02 := getPackage("The tragedy resounded around the world.")

		rNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(pkg01))

		parser(pkg01[:rNum])
		parser(append(pkg01[rNum:], pkg02...))

		So(len(results), ShouldEqual, 2)
		So(string(results[0]), ShouldEqual, "The news rocked the world.")
		So(string(results[1]), ShouldEqual, "The tragedy resounded around the world.")
	})

	Convey("两个包，后面一个包被分包，前半部分和前面的粘包", t, func() {
		var results [][]byte

		parser := ParserFactory(func(b []byte) {
			results = append(results, b)
		})

		pkg01 := getPackage("The news rocked the world.")
		pkg02 := getPackage("The tragedy resounded around the world.")

		rNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(pkg02))

		parser(append(pkg01, pkg02[:rNum]...))
		parser(pkg02[rNum:])

		So(len(results), ShouldEqual, 2)
		So(string(results[0]), ShouldEqual, "The news rocked the world.")
		So(string(results[1]), ShouldEqual, "The tragedy resounded around the world.")
	})
}
