package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func checkError1(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func TestClient(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "2.gzzhwy.cc:7100")
	checkError1(err)

	//建立tcp连接
	conn, err := net.DialTCP("tcp", nil, addr)
	checkError1(err)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	//向服务端发送数据
	go func() {
		index := 0
		for {
			_, err = conn.Write([]byte(strconv.Itoa(index) + "号数据"))
			if err != nil {
				break
			}
			time.Sleep(time.Second)
			checkError1(err)
			index++
		}
		wg.Done()
	}()
	go func() {
		for {
			bytes := make([]byte, 1024)
			readLen, err := conn.Read(bytes)
			if err != nil {
				break
			}
			fmt.Println("收到: ", readLen, err, string(bytes[:readLen]))
		}
		wg.Done()
	}()
	wg.Wait()
}
