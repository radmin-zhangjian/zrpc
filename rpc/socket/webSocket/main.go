package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读的缓冲大小
	WriteBufferSize: 1024, // 写的缓冲大小
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data struct {
	Ip      string `json:"ip"`
	Message string `json:"message"`
	User    string `json:"user"`
	Type    string `json:"type"`
}

type connection struct {
	ws     *websocket.Conn
	output chan []byte
	data   *Data
}

type hub struct {
	mu sync.Mutex
	// connection 连接
	connections map[*connection]bool
	// 从服务器发送的信息
	broadcast chan []byte
	// 从连接器注册请求
	register chan *connection
	// 销毁请求
	unregister chan *connection
}

var h = &hub{
	connections: make(map[*connection]bool),
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

func (h *hub) Run() {
	for {
		select {
		case wsc := <-h.register:
			h.mu.Lock()
			h.connections[wsc] = true
			h.mu.Unlock()
			// 组装data数据
			wsc.data.Ip = wsc.ws.RemoteAddr().String()
			// 更新类型
			wsc.data.Type = "handshake"
			wsc.data.User = ""
			wsc.data.Message = string("hello")
			dataByte, _ := json.Marshal(wsc.data)
			wsc.output <- dataByte
		case wsc := <-h.unregister:
			// 判断map里是否存在要删的数据
			if _, ok := h.connections[wsc]; ok {
				h.mu.Lock()
				delete(h.connections, wsc)
				h.mu.Unlock()
				// 关闭连接管道
				close(wsc.output)
			}
		case data := <-h.broadcast:
			h.Connections(data)
		}
	}
}

func (h *hub) Connections(data []byte) {
	// 广播所有人
	for c := range h.connections {
		select {
		case c.output <- data:
		default:
			// 防止死循环
			mu := new(sync.Mutex)
			mu.Lock()
			delete(h.connections, c)
			mu.Unlock()
			close(c.output)
		}
	}
}

func (wsc *connection) Writer() {
	// 写入ws数据
	for message := range wsc.output {
		err := wsc.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			continue
		}
	}
	wsc.ws.Close()
}

func (wsc *connection) Reader() {
	for {
		// 读取ws中的数据
		_, message, err := wsc.ws.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "ping" {
			message = []byte("pong")
		}

		log.Println("message: ", string(message))
		json.Unmarshal(message, &wsc.data)
		//log.Printf("message: %#v", wsc.data)

		switch wsc.data.Type {
		case "ping":
			wsc.data.Type = "pong"
			dataByte, _ := json.Marshal(wsc.data)
			h.broadcast <- dataByte
		case "login":
			wsc.data.Type = "login"
			dataByte, _ := json.Marshal(wsc.data)
			h.broadcast <- dataByte
		case "user":
			wsc.data.Type = "user"
			data_byte, _ := json.Marshal(wsc.data)
			h.broadcast <- data_byte
		case "logout":
			wsc.data.Type = "logout"
			data_byte, _ := json.Marshal(wsc.data)
			h.broadcast <- data_byte
			h.unregister <- wsc
		}
	}
}

func ping(c *gin.Context) {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	wsc := &connection{
		ws:     ws,
		output: make(chan []byte, 128),
		data:   new(Data),
	}

	// 注册
	h.register <- wsc

	go wsc.Writer()
	wsc.Reader()

	defer func() {
		wsc.data.Type = "logout"
		data_byte, _ := json.Marshal(wsc.data)
		h.broadcast <- data_byte
		h.unregister <- wsc
	}()
}

func longconnect(c *gin.Context) {
	c.HTML(http.StatusOK, "longconnect.html", gin.H{})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("rpc/socket/webSocket/templates/*")
	// Ping test
	r.GET("/ping", ping)
	r.GET("/longconnect.html", longconnect)
	return r
}

func main() {
	go h.Run()
	r := setupRouter()
	r.Run(":9090")
}
