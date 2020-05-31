package main

import (
	"GowGateServer/tools"
	"fmt"
	"net"
	"sync"
)

//收集GameGate连接DB时数据交互
func main() {
	tcpServer, _ := net.ResolveTCPAddr("tcp", ":6000")
	listener, _ := net.ListenTCP("tcp", tcpServer)
	fmt.Println("开始监听6000端口的DB连接...")
	userConn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("6000端口的DB连接成功")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	handleUserConnection(userConn, connectGowServer(), wg)
	wg.Wait()
}

//连接战神GateGameServer
func connectGowServer() *net.TCPConn {
	fmt.Println("开始连接6001端口的M2服务...")
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6001")
	tools.CheckError(err)
	//建立tcp连接
	conn, err := net.DialTCP("tcp", nil, addr)
	tools.CheckError(err)
	fmt.Println("连接战神M2服务(6001)成功")
	return conn
}

func handleUserConnection(ggs net.Conn, db *net.TCPConn, wg *sync.WaitGroup) {
	g2u := &tools.ProxyChannel{SourceConn: db, TargetConn: ggs, Flag: tools.SERVER_2_USER}
	u2g := &tools.ProxyChannel{SourceConn: ggs, TargetConn: db, Flag: tools.USER_2_SERVER}
	logChannel := make(chan string, 2)
	g2u.LogChannel = logChannel
	u2g.LogChannel = logChannel
	go tools.ChanelProxyData(g2u)
	go tools.ChanelProxyData(u2g)
	tools.NewLog("DB-M2交互日志.txt").LoopWriteLog(logChannel, wg, nil)
}
