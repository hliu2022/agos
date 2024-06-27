package log

import (
	"agos-server/libs/agos/core/buffer"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var LogSize = int64(40960)
var OpenConsole = true
var LogChannelCap = 4096
var LogPath = "./log/"
var LogLevel = LevelDebug
var sLog Logger

type ELevelType int

// levels
const (
	LevelDebug ELevelType = 0
	LevelInfo  ELevelType = 1
	LevelWarn  ELevelType = 2
	LevelError ELevelType = 3
	LevelStack ELevelType = 4
	LevelMax
)

var colorStr = [LevelMax + 1]string{
	"\033[1;35m ",
	"\033[32m ",
	"\033[1;33m ",
	"\033[1;4;31m ",
	"\033[1;4;31m ",
}

var formatStr = [LevelMax + 1]string{
	"[DEBUG]",
	"[ INFO]",
	"[ WARN]",
	"[ERROR]",
	"[ERROR]",
}

func InitLogger(level ELevelType, logPath string, filePrefix string, logChannelCap int) error {
	LogLevel = level
	LogPath = logPath
	sLog.writer.filePath = logPath
	sLog.writer.filePrefix = filePrefix

	sLog.writer.setLogChannel(logChannelCap)
	err := sLog.writer.switchLogFile()
	if err != nil {
		return err
	}
	return nil
}

type Logger struct {
	bf     buffer.Buffer
	writer IOWriter
}

func (s *Logger) writeLog(level ELevelType, v ...interface{}) {
	if level < LogLevel {
		return
	}
	s.bf.Reset()

	s.bf.AppendString(colorStr[level])
	s.formatHeader(&s.bf, level, 3, v...)
	s.bf.AppendString("\u001B[0m")

	//for _, sr := range a {
	//	s.bf.AppendString(slog.AnyValue(sr).String())
	//}
	s.bf.AppendString("\n")
	s.writer.Write([]byte(s.bf.Bytes()))
}

func (s *Logger) formatHeader(buf *buffer.Buffer, level ELevelType, depth int, v ...interface{}) {
	now := time.Now()
	var file string
	var line int

	// Release lock while getting caller info - it's expensive.
	var ok bool
	_, file, line, ok = runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)

	buf.AppendString(
		fmt.Sprintf("%s %s %s %v",
			now.Format("2006/01/02 15:04:05"),
			formatStr[level],
			strings.Join([]string{file, strconv.Itoa(line)}, ":"),
			fmt.Sprint(v...),
		),
	)

	if level == LevelStack {
		buf.AppendString("\n")
		buf.AppendBytes(debug.Stack())
	}
}

func Debug(v ...interface{}) {
	sLog.writeLog(LevelDebug, v...)
}

func Info(v ...interface{}) {
	sLog.writeLog(LevelInfo, v...)
}

func Warn(v ...interface{}) {
	sLog.writeLog(LevelWarn, v...)
}

func Error(v ...interface{}) {
	sLog.writeLog(LevelError, v...)
}

func Stack(v ...interface{}) {
	sLog.writeLog(LevelStack, v...)
}
