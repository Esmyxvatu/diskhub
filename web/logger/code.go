package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct{}

//============================================================ Functions ============================================================================

func (l *Logger) Info(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	_, filepath, line, _ := runtime.Caller(1)
	file := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileName := strings.Split(file, ".")[0]

	level := "\033[42m INFO  \033[0m"

	fmt.Printf("\033[1m%s\033[0m │%s│ %-20s │ %s\n",
		time.Now().Format("2006/01/02 15:04:05.000"), // Date/Hour
		level,                         // Level
		fileName+":"+fmt.Sprint(line), // File:Line
		logMessage,                    // Message
	)
}
func (l *Logger) Error(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	_, filepath, line, _ := runtime.Caller(1)
	file := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileName := strings.Split(file, ".")[0]

	level := "\033[45m ERROR \033[0m"

	fmt.Printf("\033[1m%s\033[0m │%s│ %-20s │ %s\n",
		time.Now().Format("2006/01/02 15:04:05.000"), // Date/Hour
		level,                         // Level
		fileName+":"+fmt.Sprint(line), // File:Line
		logMessage,                    // Message
	)
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	_, filepath, line, _ := runtime.Caller(1)
	file := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileName := strings.Split(file, ".")[0]

	level := "\033[41m FATAL \033[0m"

	fmt.Printf("\033[1m%s\033[0m │%s│ %-20s │ %s\n",
		time.Now().Format("2006/01/02 15:04:05.000"), // Date/Hour
		level,                         // Level
		fileName+":"+fmt.Sprint(line), // File:Line
		logMessage,                    // Message
	)

	os.Exit(1)
}
func (l *Logger) Verify(err error) {
	if err != nil {
		l.Error("%s", err.Error())
	}
}
