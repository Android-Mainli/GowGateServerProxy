package main

import (
	"fmt"
	"io"
	"net"
)

//模拟server端
func main() {
	tcpServer, _ := net.ResolveTCPAddr("tcp", ":7101")
	listener, _ := net.ListenTCP("tcp", tcpServer)
	for {
		//当有新的客户端请求来的时候，拿到与客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("收到用户连接----------------------------")
		//处理逻辑
		go handle(conn)
	}
}

func handle(userConnect net.Conn) {
	userConnect.Write([]byte("收到你的连接,开始游戏吧"))
	cp := &ProxyWrite{userConnect}
	fmt.Println("收到你的连接,开始游戏吧------------")
	for {
		_, err := io.CopyBuffer(cp, cp, nil)
		if err != nil {
			cp.Close()
			break
		}
	}
}

type ProxyWrite struct {
	base net.Conn
}

func (cp *ProxyWrite) Close() error {
	return cp.base.Close()
}

// LocalAddr returns the local network address.
func (cp *ProxyWrite) LocalAddr() net.Addr {
	return cp.base.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (cp *ProxyWrite) RemoteAddr() net.Addr {
	return cp.base.RemoteAddr()
}

func (cp *ProxyWrite) Read(b []byte) (n int, err error) {
	read, err := cp.base.Read(b)
	if err == nil {
		fmt.Println("GOW读取数据", string(b[:read]))
	}
	return read, err
}

func (cp *ProxyWrite) Write(b []byte) (n int, err error) {
	fmt.Println("GOW写入数据", string(b[:len(b)]))
	write, err := cp.base.Write(b)
	return write, err
}
