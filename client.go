package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn

	return client
}

// DealResponse 处理 server 回应的消息，直接显示到标准输出即可
func (c *Client) DealResponse() {
	// 一旦 c.conn 有数据，就直接 copy 到 stdout 标准输出上，永久阻塞监听
	io.Copy(os.Stdout, c.conn)
}

// menu 菜单选择
func (c *Client) menu() bool {
	var f int

	fmt.Println(">>>>>> 1. 公聊模式")
	fmt.Println(">>>>>> 2. 私聊模式")
	fmt.Println(">>>>>> 3. 更新用户名")
	fmt.Println(">>>>>> 0. 退出")

	fmt.Scanln(&f)

	if f >= 0 && f <= 3 {
		c.flag = f
		return true
	} else {
		fmt.Println(">>>>>> 请输入合法范围内的数字 <<<<<<<")
		return false
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>> 请输入用户名：")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		// 根据不同的模式处理不同的业务
		switch c.flag {
		case 1:
			// 公聊模式
			break
		case 2:
			// 私聊模式
			break
		case 3:
			// 更新用户名
			c.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// init 初始化命令行参数，使用 ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器 IP 地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> 连接服务器失败...")
		return
	}

	// 单独开启一个 goroutine 去处理 server 回执的消息
	go client.DealResponse()

	fmt.Println(">>>>>> 连接服务器成功...")

	client.Run()
}
