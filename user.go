package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	Message chan string
	conn net.Conn
	server *Server
}

//监听消息，一旦有消息，则传给客户端
func (u *User)listenMessage()  {
	for {
		msg:= <- u.Message
		u.conn.Write([]byte(msg + "\n"))
	}
}

func (u *User) SendMsg(msg string)  {
	u.conn.Write([]byte(msg))
}

func (u *User) HandleMessage(msg string) {
	if msg == "who" {//返回当前在线用户
		u.server.Locker.Lock()
		for _, user := range u.server.OnlineUsers {
			msg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(msg)
		}
		u.server.Locker.Unlock()
	} else if len(msg) > 7 && msg[:7]=="rename|"  {//重命名，格式为"rename|jack"
		name := msg[7:]
		u.server.Locker.Lock()
		_, ok := u.server.OnlineUsers[name]
		if ok {
			u.SendMsg("用户名已存在")
		} else {
			delete(u.server.OnlineUsers, name)
			u.server.OnlineUsers[name] = u
			u.Name = name
			u.SendMsg("用户名已修改为：" + u.Name )
		}
		u.server.Locker.Unlock()

	} else if len(msg) > 4 && msg[:3]=="to|" {//私聊，格式为"to|jack|hello"
		strArr := strings.Split(msg, "|")
		if len(strArr) != 3 {
			u.SendMsg("私聊格式为\"to|jack|hello\"\n")
		} else {
			toWho := strArr[1]
			msg := strArr[2]
			if who := u.server.OnlineUsers[toWho]; who != nil {
				who.SendMsg(msg + "\n")
			} else {
				u.SendMsg(toWho + "未上线")
			}
		}
	} else {
		u.server.broadcastMessage(u, msg)
	}
}
func (u *User) Online()  {
	u.server.Locker.Lock()
	u.server.OnlineUsers[u.Name] = u
	u.server.Locker.Unlock()
	u.server.broadcastMessage(u, "已上线")
}
func (u *User) Offline() {
	u.server.Locker.Lock()
	delete(u.server.OnlineUsers, u.Name)
	u.server.Locker.Unlock()
	u.server.broadcastMessage(u, "下线")
}
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{addr, addr, make(chan string), conn, server}
	go user.listenMessage()
	return user
}