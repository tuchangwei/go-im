package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	conn net.Conn
	Name string
	menuType int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIp: serverIP,
		ServerPort: serverPort,
		menuType: 999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		return nil
	}
	client.conn = conn
	return client
}
func (c *Client) menu() bool {
	var menuType int
	fmt.Scanln(&menuType)
	if menuType >= 0 && menuType <= 3 {
		c.menuType = menuType
		return true
	}
	fmt.Println("请输入0到3的合法字符")
	return false
}
func (c *Client) rename() {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&c.Name)
	msg := "rename|" + c.Name + "\n"
 	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("c.conn.Write error:", err)
	}
}
func (c *Client) PublicChat() {
	fmt.Println("请输入消息，exit表示退出")
	for {
		var msg string
		fmt.Scanln(&msg)
		fmt.Println(">>>>" + msg)
		if msg == "exit" {
			fmt.Println("已退出公聊模式")
			break
		}
		if msg != "" {
			_, err := c.conn.Write([]byte(msg+"\n"))
			if err != nil {
				fmt.Println("c.conn.Write error:", err)
				break
			}
		}
	}
}
func (c *Client) PrivateChat() {
	for {
		fmt.Println("请输入用户名，exit表示退出")
		_, err := c.conn.Write([]byte("who"+"\n"))
		if err != nil {
			return
		}
		var name string
		fmt.Scanln(&name)
		if name == "exit" {
			return
		}
		fmt.Printf("开始和%s聊天吧，请输入聊天内容, exit表示退出聊天\n", name)
		for {
			var msg string
			fmt.Scanln(&msg)
			if msg == "exit" {
				break
			}
			if msg != "" {
				_, err := c.conn.Write([]byte("to|"+ name + "|" + msg+"\n"))
				if err != nil {
					fmt.Println("c.conn.Write error:", err)
				}
			}

		}




	}

}

func (c *Client) HandleResponse()  {
	io.Copy(os.Stdout, c.conn)
}
func (c *Client) Run() {
	for  {
		fmt.Println("请选择模式：")
		fmt.Println("1,公聊模式")
		fmt.Println("2,私聊模式")
		fmt.Println("3,重命名")
		fmt.Println("0,退出")
		if c.menu() {
			switch c.menuType {
			case 0://退出
				fmt.Println("退出")
				return
			case 1://公聊模式
				fmt.Println("公聊模式")
				c.PublicChat()
				break
			case 2://私聊模式
				fmt.Println("私聊模式")
				c.PrivateChat()
				break
			case 3://重命名
				fmt.Println("重命名模式")
				c.rename()
				break
			}
		}
	}
}
var serverIP string
var serverPort int
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器IP（默认值为：127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8787, "设置服务器Port（默认值为：8787)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")
	go client.HandleResponse()
	client.Run()
}