package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser 创建一个用户的 API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前 User channel 的方法，一旦有消息，就直接发送给对端客户端
func (s *User) ListenMessage() {
	for {
		msg := <-s.C

		s.conn.Write([]byte(msg + "\n"))
	}
}
