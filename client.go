package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

// DealResponse 处理 server 回应的消息，直接显示到终端即可
func (c *Client) DealResponse() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from server: %v", err)
		}

		if strings.TrimSpace(message) == "你被踢了" {
			fmt.Println("你被踢了")
			os.Exit(1)
		}

		fmt.Print(message)
	}
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

// SelectUsers 查询在线用户
func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// PrivateChat 私聊模式
func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	fmt.Println(">>>>>> 请输入聊天对象用户名，exit 退出：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		if remoteName == c.Name {
			fmt.Println(">>>>>> 不能和自己聊天，请重新输入：")
			fmt.Scanln(&remoteName)
			continue
		}

		fmt.Println(">>>>>> 请输入聊天内容，exit 退出：")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}

			chatMsg = ""
			fmt.Println(">>>>>> 请输入聊天内容，exit 退出：")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUsers()
		fmt.Println(">>>>>> 请输入聊天对象用户名，exit 退出：")
		fmt.Scanln(&remoteName)
	}
}

// PublicChat 公聊模式
func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>> 请输入聊天内容，exit 退出：")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		_, err := c.conn.Write([]byte(chatMsg + "\n"))
		if err != nil {
			fmt.Println("conn.Write err:", err)
			break
		}

		chatMsg = ""
		fmt.Println(">>>>>> 请输入聊天内容，exit 退出：")
		fmt.Scanln(&chatMsg)
	}
}

// UpdateName 更新用户名
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
			c.PublicChat() // 公聊模式
			break
		case 2:
			c.PrivateChat() // 私聊模式
			break
		case 3:
			c.UpdateName() // 更新用户名
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
