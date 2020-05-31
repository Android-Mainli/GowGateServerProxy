package tools

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	USER_2_SERVER = "USER - SERVER"
	SERVER_2_USER = "SERVER - USER"
)

type ProxyChannel struct {
	SourceConn net.Conn
	TargetConn net.Conn
	Flag       string
	LogChannel chan string
	last       int64
}

func (cp *ProxyChannel) Close() error {
	return cp.SourceConn.Close()
}

// LocalAddr returns the local network address.
func (cp *ProxyChannel) LocalAddr() net.Addr {
	return cp.SourceConn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (cp *ProxyChannel) RemoteAddr() net.Addr {
	return cp.SourceConn.RemoteAddr()
}

func (cp *ProxyChannel) Read(b []byte) (n int, err error) {
	read, err := cp.SourceConn.Read(b)
	if err == nil {
		if cp.Flag == USER_2_SERVER {
			cp.log(b, read, "")
		}
		//fmt.Println(cp.Flag, "读取数据", b[:read])
	}
	return read, err
}

func (cp *ProxyChannel) Write(b []byte) (n int, err error) {
	if cp.Flag == SERVER_2_USER {
		cp.log(b, len(b), "")
		//fmt.Println(cp.Flag, "写入数据", b[:len(b)])
	}
	write, err := cp.TargetConn.Write(b)
	return write, err
}

func (cp *ProxyChannel) log(b []byte, read int, hz string) {
	if cp.last > 0 {
		cp.last = time.Now().Unix() - cp.last
	}
	cp.LogChannel <- fmt.Sprintf("\r\n\n\n%d - %s:\n\n%v\n\n%s\n\n", cp.last, cp.Flag, b[:read], string(b[:read]))
	cp.last = time.Now().Unix()
}

func ChanelProxyData(channel *ProxyChannel) {
	_, err := io.CopyBuffer(channel, channel, nil)
	if err != nil {
		channel.Close()
	}
}
