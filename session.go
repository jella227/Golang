/********************************************************/
// WebSocket Session对象
// Author 		:Jella
// Version 		:1.0.2(release)
// Dependency	:github.com/gorilla/websocket
/********************************************************/

package ws

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

//写线程
func (session *Session) write() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-session.out:
		case <-session.cls:
			goto ERR
		}
		if err = session.ws.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	session.Close()
}

//读线程（由Session对象内部线程操作）
func (session *Session) read() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = session.ws.ReadMessage(); err != nil {
			goto ERR
		}
		if len(data) > config.BufferLen {
			fmt.Println("[Error]: client send data length overflow. data length > " + strconv.Itoa(config.BufferLen) + "byte.")
			continue
		}
		select {
		case session.in <- data:
		case <-session.cls:
			goto ERR
		}
	}
ERR:
	session.Close()
}

/**
 * 读取消息（由服务线程进行操作）
 * @param handler 消息处理函数
 */
func (session *Session) reciMessage(handler messageHandler) {
	var (
		data     []byte
		callback messageHandler = handler
	)
	for {

		select {
		case data = <-session.in:
			if callback != nil {
				callback(session, data)
			}

		case <-session.cls:
			goto ERR
		}

	}
ERR:
	session.Close()
}

/**
 * 创建连接对象
 * @param wsc websocket连接对象
 * @param len 接收与发送数据的长度（单位：字节）
 * @return 连接对象，错误信息
 */
func createWsSession(wsc *websocket.Conn, len int) (session *Session, err error) {
	session = &Session{
		ws:            wsc,
		RemoteAddress: wsc.RemoteAddr().String(),
		in:            make(chan []byte, len),
		out:           make(chan []byte, len),
		cls:           make(chan byte, 1),
	}

	go session.read()  //读线程
	go session.write() //写线程

	return
}

///////////////////////////////////////////////////////
/////////////////////// 外部接口 ///////////////////////
///////////////////////////////////////////////////////

/**
 * 连接对象结构体（应由服务来创建，外部不可手动创建）
 */
type Session struct {
	//websocket相关
	ws            *websocket.Conn
	RemoteAddress string

	//读、写相关
	in  chan []byte
	out chan []byte

	//关闭连接相关
	cls     chan byte
	mutex   sync.Mutex
	isClose bool

	//其他
	Params interface{}
}

/**
 * 关闭Session连接
 *
 */
func (session *Session) Close() {
	session.ws.Close() //线程安全，可多次调用
	session.mutex.Lock()
	if !session.isClose {
		close(session.cls) //用于关闭close channel通道对象
		session.isClose = true

		//触发断开连接后续事件
		sessionBreak(session)
	}
	session.mutex.Unlock()
}

/**
 * 连接是否处于关闭状态
 * @return true：关闭；false：打开
 */
func (session *Session) IsClosed() bool {
	if session != nil {
		return session.isClose
	}
	return false
}

/**
 * 发送消息
 * @param data 数据
 */
func (session *Session) SendMessage(data []byte) {
	if len(data) > config.BufferLen {
		fmt.Println("[Error]: server send data length overflow. data length > " + strconv.Itoa(config.BufferLen) + "byte.")
		return
	}

	select {
	case session.out <- data:
	case <-session.cls:
		session.Close()
	}
}
