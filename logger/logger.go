package logger

import (
	"fmt"
	"os"
)

const (
	FATAL = 0
	PANIC = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
)

var Level int64 = INFO

var lg *Logger

func init() {
	lg = New(os.Stdout, "")
}

func Info(format string, v ...interface{}) {
	if Level >= INFO {
		lg.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...))
	}
}

func Warn(format string, v ...interface{}) {
	if Level >= WARN {
		lg.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v...))
	}
}

func Error(format string, v ...interface{}) {
	if Level >= ERROR {
		msg := fmt.Sprintf("[ERROR] "+format+"\n", v...)
		lg.Output(2, msg)
	}
}

func ErrorN(n int, format string, v ...interface{}) {
	if Level >= ERROR {
		msg := fmt.Sprintf("[ERROR] "+format+"\n", v...)
		lg.Output(2+n, msg)
	}
}

func Debug(format string, v ...interface{}) {
	if Level >= DEBUG {
		lg.Output(2, fmt.Sprintf("[DEBUG] "+format+"\n", v...))
	}
}

func Fatal(format string, v ...interface{}) {
	if Level >= FATAL {
		msg := fmt.Sprintf("[FATAL] "+format+"\n", v...)
		lg.Output(2, msg)
		os.Exit(1)
	}
}

func Panic(format string, v ...interface{}) {
	if Level >= PANIC {
		msg := fmt.Sprintf("[PANIC] "+format+"\n", v...)
		lg.Output(2, msg)
		panic(msg)
	}
}