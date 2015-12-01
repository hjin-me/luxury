package logex

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	NONE Level = iota
	FATAL
	WARNING
	NOTICE
	TRACE
	DEBUG
	LEVEL_MAX
)

type Level uint
type Logger interface {
	Id() string
	SetColor(bool)
	SetLogLevel(Level)
	SetOutput(io.Writer)
	Output(Level, int, string) error
	Debug(...interface{})
	Debugf(string, ...interface{})
	Trace(...interface{})
	Tracef(string, ...interface{})
	Notice(...interface{})
	Noticef(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

type baseLogger struct {
	id      string
	w       io.Writer
	buf     []byte
	mu      sync.Mutex
	enabled [LEVEL_MAX]bool
	color   bool
}

var levelStr = []string{"NONE", "FATAL", "WARNING", "NOTICE", "TRACE", "DEBUG"}

var (
	std = newRaw(DEBUG, os.Stdout, true)
)

func LogId() string {
	var reqId = time.Now().UnixNano()
	return strconv.FormatInt(reqId, 10)
}

/*
func NewWithLogId(logId string) *log.Logger {
	w := getFD()
	return log.New(w, logId+": ", log.Ldate|log.Ltime|log.Llongfile)
}
*/
func New() Logger {
	return &baseLogger{id: LogId()}
}

func newRaw(level Level, w io.Writer, color bool) Logger {
	std := &baseLogger{id: LogId()}
	std.SetLogLevel(level)
	std.SetOutput(w)
	std.SetColor(color)
	return std
}

/*
func getFD() io.Writer {
	once.Do(func() {
		cfg := Config{}
		err := config.Load("log.yaml", &cfg)
		log.Println(cfg)
		if err != nil {
			log.Println("can not get conf", "use stdout as log output")
			fd = os.Stdout
			return
		}
		fd, err = os.Open(cfg.Dest)
		if err != nil {
			log.Println("can not open file ", cfg.Dest, "use stdout as log output")
			fd = os.Stdout
			return
		}
	})
	fd = os.Stdout
	return fd
}
*/

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func SetLogLevel(l Level) {
	std.SetLogLevel(l)
}
func SetColor(b bool) {
	std.SetColor(b)
}

func (l *baseLogger) SetLogLevel(level Level) {
	for i := FATAL; i < LEVEL_MAX; i++ {
		if i <= level {
			l.enabled[i] = true
		} else {
			l.enabled[i] = false
		}
	}
}

func (l *baseLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.w = w
}

func (l *baseLogger) SetColor(b bool) {
	l.color = b
}

func (l *baseLogger) Id() string {
	return l.id
}

func (l *baseLogger) Output(level Level, callDepth int, s string) error {
	if level > DEBUG {
		return errors.New("wrong log level")
	} else if !l.enabled[level] {
		return nil
	}

	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		file = "???"
		line = 0
	}

	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf = l.buf[:0]
	l.formatPrefix(level, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	if l.color {
		l.buf = append(l.buf, "\033[0m"...)
	}
	_, err := l.w.Write(l.buf)

	return err
}

func (l *baseLogger) formatPrefix(level Level, t time.Time, file string, line int) {
	if l.color {
		switch level {
		case FATAL:
			l.buf = append(l.buf, "\033[0;31m"...)
		case WARNING:
			l.buf = append(l.buf, "\033[0;33m"...)
		case NOTICE:
			l.buf = append(l.buf, "\033[0;32m"...)
		case TRACE:
			l.buf = append(l.buf, "\033[0;34m"...)
		case DEBUG:
			l.buf = append(l.buf, "\033[0m"...)
		}
	}
	l.buf = append(l.buf, levelStr[level]...)
	if l.color {
		l.buf = append(l.buf, "\033[0m"...)
	}
	l.buf = append(l.buf, ": "...)
	_, month, day := t.Date()
	l.buf = itoa(l.buf, int(month), 2)
	l.buf = append(l.buf, '-')
	l.buf = itoa(l.buf, day, 2)
	l.buf = append(l.buf, ' ')
	hour, min, sec := t.Clock()
	l.buf = itoa(l.buf, hour, 2)
	l.buf = append(l.buf, ':')
	l.buf = itoa(l.buf, min, 2)
	l.buf = append(l.buf, ':')
	l.buf = itoa(l.buf, sec, 2)
	l.buf = append(l.buf, '.')
	l.buf = itoa(l.buf, t.Nanosecond()/1e6, 3)
	l.buf = append(l.buf, ": "...)

	l.buf = itoa(l.buf, int(1), -1)
	l.buf = append(l.buf, ": "...)

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	l.buf = append(l.buf, short...)
	l.buf = append(l.buf, ':')
	l.buf = itoa(l.buf, line, -1)
	l.buf = append(l.buf, ": "...)
	if l.color {
		l.buf = append(l.buf, "\033[0m"...)
	}
}

func itoa(buf []byte, i int, wid int) []byte {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		buf = append(buf, '0')
		return buf
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	buf = append(buf, b[bp:]...)
	return buf
}

// Fatalf is equivalent to Printf() for FATAL-level log.
func Fatalf(format string, v ...interface{}) {
	std.Output(FATAL, 2, fmt.Sprintf(format, v...))
}

// Fatal is equivalent to Print() for FATAL-level log.
func Fatal(v ...interface{}) {
	std.Output(FATAL, 2, fmt.Sprintln(v...))
}

// Warningf is equivalent to Printf() for WARNING-level log.
func Warningf(format string, v ...interface{}) {
	std.Output(WARNING, 2, fmt.Sprintf(format, v...))
}

// Waring is equivalent to Print() for WARING-level log.
func Warning(v ...interface{}) {
	std.Output(WARNING, 2, fmt.Sprintln(v...))
}

// Noticef is equivalent to Printf() for NOTICE-level log.
func Noticef(format string, v ...interface{}) {
	std.Output(NOTICE, 2, fmt.Sprintf(format, v...))
}

// Notice is equivalent to Print() for NOTICE-level log.
func Notice(v ...interface{}) {
	std.Output(NOTICE, 2, fmt.Sprintln(v...))
}

// Tracef is equivalent to Printf() for TRACE-level log.
func Tracef(format string, v ...interface{}) {
	std.Output(TRACE, 2, fmt.Sprintf(format, v...))
}

// Trace is equivalent to Print() for TRACE-level log.
func Trace(v ...interface{}) {
	std.Output(TRACE, 2, fmt.Sprintln(v...))
}

// Debugf is equivalent to Printf() for DEBUG-level log.
func Debugf(format string, v ...interface{}) {
	std.Output(DEBUG, 2, fmt.Sprintf(format, v...))
}

// Debug is equivalent to Print() for DEBUG-level log.
func Debug(v ...interface{}) {
	std.Output(DEBUG, 2, fmt.Sprintln(v...))
}

// Fatalf is equivalent to Printf() for FATAL-level log.
func (l *baseLogger) Fatalf(format string, v ...interface{}) {
	l.Output(FATAL, 2, fmt.Sprintf(format, v...))
}

// Fatal is equivalent to Print() for FATAL-level log.
func (l *baseLogger) Fatal(v ...interface{}) {
	l.Output(FATAL, 2, fmt.Sprintln(v...))
}

// Warningf is equivalent to Printf() for WARNING-level log.
func (l *baseLogger) Warningf(format string, v ...interface{}) {
	l.Output(WARNING, 2, fmt.Sprintf(format, v...))
}

// Waring is equivalent to Print() for WARING-level log.
func (l *baseLogger) Warning(v ...interface{}) {
	l.Output(WARNING, 2, fmt.Sprintln(v...))
}

// Noticef is equivalent to Printf() for NOTICE-level log.
func (l *baseLogger) Noticef(format string, v ...interface{}) {
	l.Output(NOTICE, 2, fmt.Sprintf(format, v...))
}

// Notice is equivalent to Print() for NOTICE-level log.
func (l *baseLogger) Notice(v ...interface{}) {
	l.Output(NOTICE, 2, fmt.Sprintln(v...))
}

// Tracef is equivalent to Printf() for TRACE-level log.
func (l *baseLogger) Tracef(format string, v ...interface{}) {
	l.Output(TRACE, 2, fmt.Sprintf(format, v...))
}

// Trace is equivalent to Print() for TRACE-level log.
func (l *baseLogger) Trace(v ...interface{}) {
	l.Output(TRACE, 2, fmt.Sprintln(v...))
}

// Debugf is equivalent to Printf() for DEBUG-level log.
func (l *baseLogger) Debugf(format string, v ...interface{}) {
	l.Output(DEBUG, 2, fmt.Sprintf(format, v...))
}

// Debug is equivalent to Print() for DEBUG-level log.
func (l *baseLogger) Debug(v ...interface{}) {
	l.Output(DEBUG, 2, fmt.Sprintln(v...))
}
