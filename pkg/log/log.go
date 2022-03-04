package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// log levels
const (
	TraceLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	Disabled
)

const (
	colorRed = uint8(iota + 91)
	colorGreen
	colorYellow
	colorBlue
	colorMagenta // Red color
)

var (
	traceLog = log.New(os.Stderr, "", 0x0)
	debugLog = log.New(os.Stderr, "", 0x0)
	infoLog  = log.New(os.Stderr, "", 0x0)
	warnLog  = log.New(os.Stderr, "", 0x0)
	errorLog = log.New(os.Stderr, "", 0x0)
	fatalLog = log.New(os.Stderr, "", 0x0)
	loggers  = []*log.Logger{traceLog, errorLog, infoLog, debugLog, warnLog}

	currentLevel = Disabled
	mu           sync.Mutex
)

func red(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorRed, s)
}
func green(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorGreen, s)
}
func yellow(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorYellow, s)
}
func blue(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorBlue, s)
}
func magenta(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorMagenta, s)
}
func formatLog(prefix string) string {
	return "[" + time.Now().Format("2006/01/02 15:04:05") + "] " + prefix + " "
}

// log methods
var (
	Trace = trace
	Debug = debug
	Error = error
	Info  = info
	Warn  = warning
	Fatal = fatal // 打印错误信息，并结束进程（结束前会发送一个 CLOSE_FROM_FATAL 事件）
)

// 默认级别是 Info
func init() {
	SetLevel(InfoLevel)
}

// SetLevel 设置日志级别
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	if level > Disabled {
		level = Disabled
	}
	currentLevel = level
	// Reset output level.
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if level > TraceLevel {
		traceLog.SetOutput(ioutil.Discard)
	}

	if level > DebugLevel {
		debugLog.SetOutput(ioutil.Discard)
	}

	if level > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}

	if level > WarnLevel {
		warnLog.SetOutput(ioutil.Discard)
	}

	if level > ErrorLevel {
		errorLog.SetOutput(ioutil.Discard)
	}
}

func GetCurrentLevel() int {
	return currentLevel
}

func trace(format string, a ...interface{}) {
	traceLog.Println(formatLog(blue("[TRACE] " + fmt.Sprintf(format, a...))))
}

func info(format string, a ...interface{}) {
	infoLog.Println(formatLog(green("[INFO] " + fmt.Sprintf(format, a...))))
}

func debug(format string, a ...interface{}) {
	debugLog.Println(formatLog(yellow("[DEBUG] " + fmt.Sprintf(format, a...))))
}

func warning(format string, a ...interface{}) {
	warnLog.Println(formatLog(magenta("[WARN] " + fmt.Sprintf(format, a...))))
}

func error(format string, a ...interface{}) {
	errorLog.Println(formatLog(red("[ERROR] " + fmt.Sprintf(format, a...))))
}

func fatal(format string, a ...interface{}) {
	fatalLog.Println(formatLog(red("[FATAL] " + fmt.Sprintf(format, a...))))
	os.Exit(1)
}
