package main

import "net"

type User struct {
	Name string
	Addr string
	Message chan string
	conn net.Conn
}
//监听消息，一旦有消息，则传给客户端
func (u *User)listenMessage()  {
	for {
		msg:= <- u.Message
		u.conn.Write([]byte(msg + "\n"))
	}
}
func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{addr, addr, make(chan string), conn}
	go user.listenMessage()
	return user
}