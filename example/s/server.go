package main

import (
	"fmt"
	"Golang-master/byt"
	"Golang-master/ws"
	"strings"
)

func main() {

	var config ws.WSConfig = ws.WSConfig{
		Host:      "192.168.0.102",
		Port:      1201,
		Pattern:   "/",
		BufferLen: 4096,
	}
	ws.Listen(config, msgHandler)

}

func msgHandler(s *ws.Session, b []byte) {
	if s.IsClosed() {
		return
	}
	buf := byt.NewBufferWithByte(b)
	str := buf.ReadUTF8String()
	buf.Kill()
	buf = nil

	sary := strings.Split(str, "|")
	if len(sary) < 2 {
		return
	}
	fmt.Println(sary[0] + " -> " + sary[1])

	if !s.IsClosed() {
		_buf := byt.NewBuffer()
		_buf.WriteUTF8String("hello client!|" + sary[1])
		s.SendMessage(_buf.GetByte())
		_buf.Kill()
		_buf = nil
	}
}
