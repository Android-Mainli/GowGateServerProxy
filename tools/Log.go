package tools

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type GowLog struct {
	out *os.File
}

func NewLog(saveFileName string) *GowLog {
	file, err := ObtainNewFile(saveFileName)
	CheckError(err)
	return &GowLog{file}
}

func (g *GowLog) LoopWriteLog(logChannel chan string, group *sync.WaitGroup, ctx context.Context, proxyLog func(log string)) {
	go func(group *sync.WaitGroup, proxyLog func(log string)) {
		if group != nil {
			group.Add(1)
		}
		for {
			if log, ok := <-logChannel; ok {
				fmt.Println(log)
				if proxyLog != nil {
					proxyLog(log)
				}
				g.out.WriteString(log)
			} else {
				break
			}
		}
		if group != nil {
			group.Done()
		}
	}(group, proxyLog)
}

//获取全新最佳读写View
func ObtainNewFile(saveFileName string) (*os.File, error) {
	_, err := os.Stat(saveFileName)
	var file *os.File
	if os.IsExist(err) {
		os.Remove(saveFileName)
		file, err = os.Create(saveFileName)
	} else {
		file, err = os.Create(saveFileName)
	}
	return file, err
}
