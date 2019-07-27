package main

import (
	"fmt"
	"jella/byt"
	"jella/ws"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func main() {

	_s := make(chan os.Signal)
	signal.Notify(_s, os.Interrupt)

	//创建websocket client连接配置及连接对象
	var conf ws.WSClient_CONFIG = ws.WSClient_CONFIG{
		Host:      "192.168.0.102",
		Port:      1201,
		Path:      "/",
		BufferLen: 4096,
	}
	var cli ws.WSClient

	//注册连接关闭回调
	cli.OnClose(func() {
		fmt.Println("connect close.")
		_s <- os.Interrupt
	})

	//连接服务器
	if err := cli.Connect(conf); err != nil {
		fmt.Println("connect server fail!", err)

	} else {
		fmt.Println("connect server success!")

		//开启一条新的线程进行消息接收
		go cli.Reci(func(data []byte) {
			_buf := byt.NewBufferWithByte(data)
			sary := strings.Split(_buf.ReadUTF8String(), "|")
			fmt.Println(sary[0] + " -> " + sary[1])
			_buf.Kill()
			_buf = nil
		})

		//开启一条心的线程进行消息发送
		buf := byt.NewBuffer()
		rand.Seed(time.Now().Unix())
		numstr := ""
		go func() {
			for cli.IsConnect {
				time.Sleep(time.Millisecond * 300)
				numstr = strconv.Itoa(rand.Intn(100000000000))
				numstr += " - " + numstr

				buf.Zero()
				buf.WriteUTF8String("hello server!|" + numstr)
				cli.Send(buf.GetByte())
				break
			}
		}()
	}

	c := <-_s
	fmt.Println(c, "over!")
}
