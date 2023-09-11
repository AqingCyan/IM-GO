package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的 channel
	Message chan string
}

// ListenMessage 监听 Message 广播消息 channel 的 goroutine，一旦有消息就发送给全部的在线 User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		// 将 msg 发送给全部的在线 User
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// NewServer 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// BroadCast 广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// Handler 处理 conn 业务的接口
func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn, s)
	user.Online()
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				// 当前用户下线
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				// 读取出错且不是 EOF
				fmt.Println("conn read err:", err)
				return
			}

			// 提取用户的消息（去除 '\n'）
			msg := string(buf[:n-1])

			// 用户针对 msg 进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是一个活跃的用户
			isLive <- true
		}
	}()

	// 当前 handler 阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，应该重置定时器
			// 不做任何事情，为了激活 select，更新下面的定时器
		case <-time.After(time.Second * 300):
			// 已经超时，将当前的 User 强制关闭
			user.SendMsg("你被踢了\n")
			// 销毁用的资源
			close(user.C)
			// 关闭连接
			conn.Close()
			// 退出当前的 Handler
			return
		}
	}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	fmt.Println("服务器启动成功，等待连接...")

	// close listen socket
	defer listener.Close()

	// 启动监听 Message 的 goroutine
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}
