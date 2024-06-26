package log

import (
	"fmt"
	"io"
	"sync"
)

type LogLevelType = int

// 日志等级
const (
	LEVEL_DEBUG LogLevelType = 0
)

func Log(s string) {
	fmt.Print(s)
}

type IOWriter struct {
	file       io.Writer
	console    io.Writer
	writeBytes int64
	channels   chan []byte
	wg         sync.WaitGroup
	closeSig   chan struct{}
	lockWrite  sync.Mutex

	filePath       string
	filePrefix     string
	fileDay        int
	fileCreateTime int64
}

func (io *IOWriter) Write(bs []byte) (n int, err error) {
	return
}

func (iw *IOWriter) run() {
	defer iw.wg.Done()

Loop:
	for {
		select {
		case <-iw.closeSig:
			break Loop
		case logs := <-iw.logChannel:
			iw.writeIo(logs)
			memPool.ReleaseBytes(logs)
		}
	}

	for len(iw.logChannel) > 0 {
		logs := <-iw.logChannel
		iw.writeIo(logs)
		memPool.ReleaseBytes(logs)
	}
}
