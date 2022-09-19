package utils

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

/*
	nginx 设置来源ip
	location / {
		...
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_pass http://192.168.1.2;
		...
	}
*/

// GetCurrentIP 获取 http IP
func GetCurrentIP(r http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// GetCurrentWsIP 获取 conn IP
func GetCurrentWsIP(r *websocket.Conn) string {
	ip := r.RemoteAddr().String()
	return ip
}
