package log

import (
	"agos-server/libs/agos/core/bytespool"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

var memPool = bytespool.NewMemAreaPool()

type IOWriter struct {
	fileIo      io.Writer
	consoleIo   io.Writer
	writeBytes  int64
	logChannels chan []byte
	wg          sync.WaitGroup
	closeSig    chan struct{}
	lockWrite   sync.Mutex

	filePath       string
	filePrefix     string
	fileDay        int
	fileCreateTime int64
}

func (iw *IOWriter) Close() error {
	iw.lockWrite.Lock()
	defer iw.lockWrite.Unlock()

	iw.close()

	return nil
}

func (iw *IOWriter) close() error {
	if iw.closeSig != nil {
		close(iw.closeSig)
		iw.closeSig = nil
	}
	iw.wg.Wait()

	if iw.fileIo != nil {
		err := iw.fileIo.(io.Closer).Close()
		iw.fileIo = nil
		return err
	}

	return nil
}

func (iw *IOWriter) Write(bs []byte) (n int, err error) {
	iw.lockWrite.Lock()
	defer iw.lockWrite.Unlock()

	if iw.logChannels == nil {
		return iw.writeIo(bs)
	}

	copyBuff := memPool.MakeBytes(len(bs))
	if copyBuff == nil {
		return 0, fmt.Errorf("MakeByteSlice failed")
	}
	copy(copyBuff, bs)

	iw.logChannels <- copyBuff
	return
}

func (iw *IOWriter) writeIo(p []byte) (n int, err error) {
	iw.switchLogFile()
	if iw.fileIo != nil {
		if n, err = iw.fileIo.Write(p); n > 0 {
			atomic.AddInt64(&iw.writeBytes, int64(n))
		}
	}
	if iw.consoleIo != nil {
		n, err = iw.consoleIo.Write(p)
	}

	return
}

// 跨天或者单个文件超过限制后重新生成文件
func (iw *IOWriter) switchLogFile() error {
	now := time.Now()
	if iw.fileCreateTime == now.Unix() {
		return nil
	}
	if iw.fileDay == now.Day() && iw.isFull() == false {
		return nil
	}

	if iw.filePath == "" {
		iw.consoleIo = os.Stdout
		return nil
	}

	var err error
	fileName := fmt.Sprintf("%s%d%02d%02d_%02d_%02d_%02d.log", iw.filePrefix,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	filePath := path.Join(iw.filePath, fileName)

	iw.fileIo, err = os.Create(filePath)
	if err != nil {
		return err
	}
	iw.fileDay = now.Day()
	iw.fileCreateTime = now.Unix()
	atomic.StoreInt64(&iw.writeBytes, 0)

	if OpenConsole == true {
		iw.consoleIo = os.Stdout
	}
	return nil
}

func (iw *IOWriter) isFull() bool {
	if LogSize == 0 {
		return false
	}

	return atomic.LoadInt64(&iw.writeBytes) >= LogSize
}

func (iw *IOWriter) setLogChannel(logChannelNum int) (err error) {
	iw.lockWrite.Lock()
	defer iw.lockWrite.Unlock()
	iw.close()

	if logChannelNum == 0 {
		return nil
	}

	//copy iw.logChannel
	var logInfo []byte
	logChannel := make(chan []byte, logChannelNum)
	for i := 0; i < logChannelNum && i < len(iw.logChannels); i++ {
		logInfo = <-iw.logChannels
		logChannel <- logInfo
	}
	iw.logChannels = logChannel

	iw.closeSig = make(chan struct{})
	iw.wg.Add(1)
	go iw.run()

	return nil
}

func (iw *IOWriter) run() {
	defer iw.wg.Done()

Loop:
	for {
		select {
		case <-iw.closeSig:
			break Loop
		case logs := <-iw.logChannels:
			iw.writeIo(logs)
			memPool.ReleaseBytes(logs)
		}
	}

	for len(iw.logChannels) > 0 {
		logs := <-iw.logChannels
		iw.writeIo(logs)
		memPool.ReleaseBytes(logs)
	}
}
