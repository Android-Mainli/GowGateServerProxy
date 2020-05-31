package main

import (
	"GowGateServer/tools"
	_ "GowGateServer/tools"
	"bufio"
	"context"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"net"
	"os"
	"syscall"
	"time"
)

func connDBServer(context.Context) {
	//_, cancelFunc := context.WithCancel(context)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5100")
	checkError(err)
	//建立tcp连接
	gowServer, err := net.DialTCP("tcp", nil, addr)
	checkError(err)
	fmt.Println("连接战神DB服务(5100)成功")
	start := []byte{119, 187, 170, 51, 188, 27, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0}
	gowServer.Write(start)
	start[12] = 3
	go func() { //定时发送心跳
		for range time.Tick(time.Second * 10) {
			gowServer.Write(start)
		}
	}()
	newReader := bufio.NewReaderSize(gowServer, 1024)
	go func() { //定时发送心跳
		buff := []byte{}
		newReader.Read(buff)
		fmt.Println(buff)
	}()
}

func connM2Server(context.Context) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5000")
	checkError(err)
	//建立tcp连接
	gowServer, err := net.DialTCP("tcp", nil, addr)
	checkError(err)
	fmt.Println("连接战神M2服务(5000)成功")
	start := []byte{119, 187, 170, 51, 188, 27, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0}
	gowServer.Write(start)
	start[12] = 3
	go func() { //定时发送心跳
		for range time.Tick(time.Second * 10) {
			gowServer.Write(start)
		}
	}()
	newReader := bufio.NewReaderSize(gowServer, 1024)
	go func() { //定时发送心跳
		buff := []byte{}
		newReader.Read(buff)
		fmt.Println(buff)
	}()
}

//连接战神GateGameServer
func connectGowServer(address string) *net.TCPConn {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:7101")
	checkError(err)
	//var laddr *net.TCPAddr
	//if address != "" {
	//	laddr, err = net.ResolveTCPAddr("tcp", address)
	//	checkError(err)
	//}
	//建立tcp连接
	conn, err := net.DialTCP("tcp", nil, addr)
	checkError(err)
	fmt.Println("连接战神GameGate(7101)成功")
	return conn
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleUserConnection(ctx context.Context, userConn net.Conn, ggsConn *net.TCPConn, channel chan string) {
	g2u := &tools.ProxyChannel{SourceConn: ggsConn, TargetConn: userConn, Flag: tools.SERVER_2_USER}
	u2g := &tools.ProxyChannel{SourceConn: userConn, TargetConn: ggsConn, Flag: tools.USER_2_SERVER}
	g2u.LogChannel = channel
	u2g.LogChannel = channel
	go tools.ChanelProxyData(g2u)
	go tools.ChanelProxyData(u2g)
}

func startServer() {
	context := context.Background()
	go connDBServer(context)
	go connM2Server(context)
	tcpServer, _ := net.ResolveTCPAddr("tcp", ":7100")
	listener, _ := net.ListenTCP("tcp", tcpServer)
	infoLabel.AppendText("开始监听7100连接...\r\n")
	//fmt.Println("开始监听7100连接...")
	logChannel := make(chan string, 2)
	tools.NewLog("战神防火墙日志.txt").LoopWriteLog(logChannel, nil, context, func(log string) {
		fromString, err := syscall.UTF16FromString(log)
		if err == nil {
			infoLabel.AppendText(syscall.UTF16ToString(fromString))
		}

	})
	for {
		userConn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		s := userConn.RemoteAddr().String()
		infoLabel.AppendText(s)
		infoLabel.AppendText("用户连接...\r\n")
		//fmt.Println("收到用户连接,开启携程转发数据")
		go handleUserConnection(context, userConn, connectGowServer(s), logChannel)
	}
}

type MyMainWindow struct {
	*walk.MainWindow
}

const VERSION_NAME = "兄弟-战神网管代理 v1.0版"
const logo = "res/logo.ico"

var infoLabel *walk.TextEdit

func main() {
	mw := &MyMainWindow{}
	err := MainWindow{
		MenuItems: []MenuItem{
			Action{
				Text: "关于",
				OnTriggered: func() {
					walk.MsgBox(mw, "关于", VERSION_NAME+"\n作者：Mainli", walk.MsgBoxIconQuestion)
				},
			},
		},
		AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
		Title:    VERSION_NAME,   //标题
		MinSize:  Size{Width: 150, Height: 200},
		Size:     Size{300, 400},
		Layout:   VBox{}, //样式，纵向
		Children: []Widget{ //控件组
			Label{
				Text: "hehe",
			},
			TextEdit{
				VScroll:  true,
				ReadOnly: true,
				AssignTo: &infoLabel,
			},
		},
	}.Create() //创建
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if file, error := walk.NewIconFromFile(logo); error == nil {
		mw.SetIcon(file)
	}
	//获取屏幕宽高居中显示
	mw.SetX((int(win.GetSystemMetrics(0)) - mw.Width()) / 2)
	mw.SetY((int(win.GetSystemMetrics(1)) - mw.Height()) / 2)
	go startServer()
	mw.Run() //运行
}
