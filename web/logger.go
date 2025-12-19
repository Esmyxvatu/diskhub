package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

//============================================================ Type Definitions =====================================================================

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}
type Logger struct{}

//============================================================ Functions ============================================================================

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Save the time at the starts
		start := time.Now()

		// Create a personalize ResponseWriter to get the status code
		rr := &responseRecorder{w, http.StatusOK}

		// Pass the request to the next handler
		next.ServeHTTP(rr, r)

		// Calculate the time of the request
		duration := time.Since(start)
		duration = duration.Round(time.Millisecond)
		if duration < 0 {
			duration = 0
		}

		padding := (7 - len(fmt.Sprint(rr.statusCode))) / 2
		status := fmt.Sprintf("%s%s%s",
			strings.Repeat(" ", padding),
			fmt.Sprint(rr.statusCode),
			strings.Repeat(" ", 7-len(fmt.Sprint(rr.statusCode))-padding),
		)
		status = fmt.Sprintf("%s%s%s", getStatusColor(rr.statusCode), status, "\033[0m") // Color of the HTTP status
		method := fmt.Sprintf("%s%s%s", getMethodColor(r.Method), r.Method, "\033[0m")   // Color of the method

		// Show the log in the format wanted
		fmt.Printf("\033[1m%s\033[0m │%s│ %-20s │ %s '%s' \033[2m%s\033[0m\n",
			start.Format("2006/01/02 15:04:05.000"), // Date/Hour
			status,                                  // Code HTTP
			r.RemoteAddr,                            // IP
			method,                                  // Method
			r.URL.Path,                              // Path
			duration,                                // Duration
		)
	})
}
func (rw *responseRecorder) WriteHeader(statusCode int) {
	if rw.statusCode == http.StatusOK {
		rw.statusCode = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[42m" // Vert
	case statusCode >= 300 && statusCode < 400:
		return "\033[45m" // Rose
	case statusCode >= 400 && statusCode < 500:
		return "\033[41m" // Rouge
	default:
		return "\033[43m" // Jaune
	}
}
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Bleu
	case "POST":
		return "\033[32m" // Cyan
	case "PUT":
		return "\033[33m" // Vert
	case "DELETE":
		return "\033[31m" // Rouge
	case "PATCH":
		return "\033[36m" // Magenta
	case "OPTIONS":
		return "\033[35m" // Jaune
	default:
		return "\033[37m" // Blanc
	}
}
func (l *Logger) info(format string, args ...interface{}) {
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
func (l *Logger) error(format string, args ...interface{}) {
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
func (l *Logger) fatal(format string, args ...interface{}) {
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
func (l *Logger) verify(err error) {
	if err != nil { l.fatal(err.Error()) }
}