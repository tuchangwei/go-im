package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP          string
	Port        int
	OnlineUsers map[string]*User
	Locker      sync.RWMutex
	Message     chan string
}


//处理客户端请求
func (s *Server) handleConnection(conn net.Conn) {
	user := NewUser(conn, s)
	user.Online()
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error:", err)
				return
			}
			msg := string(buf[:n-1])
			user.HandleMessage(msg)
			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
		case <- time.After(10 * time.Second):
			user.SendMsg("你被踢了")
			s.Locker.Lock()
			defer s.Locker.Unlock()
			delete(s.OnlineUsers, user.Name)
			close(user.Message)
			conn.Close()
			return
		}

	}
}
func (s *Server) broadcastMessage(user *User, msg string)  {
	s.Message<- "[" + user.Name + "]:" + msg
}
func (s *Server) listenMessage() {
	for {
		msg := <- s.Message
		s.Locker.Lock()
		for _, user := range s.OnlineUsers {
			user.Message <- msg
		}
		s.Locker.Unlock()
	}
}

//启动服务器
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err.Error())
		return
	}
	defer listener.Close()
	for {
		//保持主程序一直运行
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err.Error())
			continue
		}
		//当有用户上线时，广播消息。
		go s.listenMessage()
		//为每个客户端连接开辟协程，处理请求
		go s.handleConnection(conn)
	}
}
//创建服务器
func NewServer(ip string, port int) *Server {
	server := &Server{ip, port, make(map[string]*User), sync.RWMutex{}, make(chan string)}
	return server
}