/********************************************************/
// WebSocket Client连接对象
// Author 		:Jella
// Version 		:1.0.2(release)
// Dependency	:github.com/gorilla/websocket
/********************************************************/
package ws

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type wsClientReci func(data []byte)
type wsClientClose func()

/** 连接配置 */
type WSClient_CONFIG struct {
	Host      string
	Port      int
	Path      string
	BufferLen int
}

/** websocket客户端 */
type WSClient struct {
	cfg       WSClient_CONFIG
	conn      *websocket.Conn
	in        chan []byte
	out       chan []byte
	cls       chan byte
	mutex     sync.Mutex
	clsFunc   wsClientClose
	err       error
	IsConnect bool
}

/**
 * 连接
 * @param conf 连接配置
 * @return 连接的错误。若没有错误则返回nil
 */
func (client *WSClient) Connect(conf WSClient_CONFIG) error {
	if client.IsConnect {
		return errors.New("已连接，不可重复操作.")
	}

	client.cfg = conf
	var addr = flag.String("addr", conf.Host+":"+strconv.Itoa(conf.Port), "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: conf.Path}

	var dialer *websocket.Dialer
	client.conn, _, client.err = dialer.Dial(u.String(), nil)
	if client.err != nil {
		client.IsConnect = false
		return client.err
	}
	client.IsConnect = true

	client.in = make(chan []byte, conf.BufferLen)
	client.out = make(chan []byte, conf.BufferLen)
	client.cls = make(chan byte, 1)

	go client.sendMessage()
	go client.reciMessage()

	return nil
}

/**
 * 发送数据
 * @param data 数据内容
 */
func (client *WSClient) Send(data []byte) {
	if !client.IsConnect {
		return
	}

	if data == nil {
		return
	}
	if len(data) > client.cfg.BufferLen {
		fmt.Println("[Error]: 数据长度溢出，无法进行发送.")
		return
	}

	select {
	case client.out <- data:
	case <-client.cls:
		client.Close()
	}
}

/**
 * 收取消息
 * @param handler 收取消息的回调函数（函数应有一个参数。参数类型[]byte。）
 */
func (client *WSClient) Reci(handler wsClientReci) {
	if !client.IsConnect {
		return
	}

	var (
		data     []byte
		callback wsClientReci = handler
	)

	for {

		select {
		case data = <-client.in:
			if callback != nil {
				callback(data)
			}

		case <-client.cls:
			goto ERR
		}

	}
ERR:
	client.Close()
}

/**
 * 注册连接关闭回调
 */
func (client *WSClient) OnClose(f wsClientClose) {
	client.clsFunc = f
}

/**
 * 关闭连接
 */
func (client *WSClient) Close() {
	client.conn.Close()

	client.mutex.Lock()
	if client.IsConnect {
		client.IsConnect = false
		close(client.cls)
		// fmt.Println("关闭连接");
		if client.clsFunc != nil {
			client.clsFunc()
		}
	}
	client.mutex.Unlock()
}

//////////////////////////////////////////////////////////
//内部实现

func (client *WSClient) sendMessage() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-client.out:
		case <-client.cls:
			goto ERR
		}
		if err = client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println(err)
			goto ERR
		}
	}
ERR:
	client.Close()
}

func (client *WSClient) reciMessage() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = client.conn.ReadMessage(); err != nil {
			goto ERR
		}
		if len(data) > client.cfg.BufferLen {
			fmt.Println("[Error]: 接收客户端发送的消息异常! 数据长度溢出。")
			continue
		}

		select {
		case client.in <- data:
		case <-client.cls:
			goto ERR
		}
	}
ERR:
	client.Close()
}
