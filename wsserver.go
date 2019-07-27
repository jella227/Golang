/********************************************************/
// WebSocket Server对象
// Author 		:Jella
// Version 		:1.0.2(release)
// Dependency	:github.com/gorilla/websocket
/********************************************************/

package ws

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

//消息处理函数
type messageHandler func(session *Session, data []byte)

//内部变量
var (
	//websocket升级协议
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	//配置
	config WSConfig

	//消息处理函数
	msgHandler messageHandler

	//连接对象
	Sessions sync.Map
)

/** 监听连接对象 */
func handler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		s      *Session
		err    error
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if s, err = createWsSession(wsConn, config.BufferLen); err != nil {
		s.Close()

	} else {
		fmt.Println("session connect. address = " + s.RemoteAddress)

		//加入缓存并实时读取消息
		Sessions.Store(s.RemoteAddress, s)
		s.reciMessage(msgHandler)
	}
}

/** Session断开处理 */
func sessionBreak(s *Session) {
	addr := s.RemoteAddress

	if msgHandler != nil {
		msgHandler(s, nil)
	}

	Sessions.Delete(addr)
	fmt.Println("session break. address = " + addr)
}

/////////////////////////////////////////////////////

/**
 * WebSocket服务配置
 */
type WSConfig struct {
	Host      string
	Port      int
	Pattern   string
	BufferLen int
}

/**
 * 启动WebSocket监听服务
 * @param conf websocket服务启动配置
 * @param mhandler 消息处理函数（函数应有2个参数，参数类型分别 *ws.Session，[]byte。第1个参数是与客户端的连接对象，第2个是消息数据）
 */
func Listen(conf WSConfig, mhandler messageHandler) {
	config = conf
	msgHandler = mhandler

	http.HandleFunc(config.Pattern, handler)
	fmt.Println("启动WebSocket服务。host=" + config.Host + " / port=" + strconv.Itoa(config.Port))

	if err := http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), nil); err != nil {
		fmt.Println("[Error]: 启动服务失败." + err.Error())
	}
}

/**
 * 广播消息
 * @param data 消息
 */
func Broadcast(data []byte) {
	if data == nil || len(data) <= 0 {
		return
	}

	Sessions.Range(func(k, v interface{}) bool {
		if !v.(*Session).IsClosed() {
			v.(*Session).SendMessage(data)
		}
		return true
	})
}

/**
 * 服务主动关闭一个客户端连接对象
 * @param k session的索引值
 */
func ShutdownClient(k string) {
	s, ok := Sessions.Load(k)
	if ok {
		if !s.(*Session).IsClosed() {
			s.(*Session).Close()
		}
	}
}
