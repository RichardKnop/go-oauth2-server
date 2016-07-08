package logging

import (
	"io"
	"log"
	"os"
)

// Level type
type level int

const (
	// Bits or'ed together to control log levels
	infoLevel level = iota
	warningLevel
	errorLevel
	fatalLevel
	flag = log.Ldate | log.Ltime
)

// Log level prefix map
var prefix = []string{
	infoLevel:    "INFO: ",
	warningLevel: "WARNING: ",
	errorLevel:   "ERROR: ",
	fatalLevel:   "FATAL: ",
}

// Logger ...
type Logger struct {
	formatter Formatter
	loggers   map[level]*log.Logger
}

// New returns instance of Logger
func New(out, errOut io.Writer, f Formatter) *Logger {
	// Fall back to stdout if out not set
	if out == nil {
		out = os.Stdout
	}

	// Fall back to stderr if errOut not set
	if errOut == nil {
		errOut = os.Stderr
	}

	// Fall back to DefaultFormatter if f not set
	if f == nil {
		f = new(DefaultFormatter)
	}

	// Initialise level loggers
	loggers := map[level]*log.Logger{
		infoLevel:    log.New(out, f.GetPrefix(infoLevel)+prefix[infoLevel], flag),
		warningLevel: log.New(out, f.GetPrefix(warningLevel)+prefix[warningLevel], flag),
		errorLevel:   log.New(errOut, f.GetPrefix(errorLevel)+prefix[errorLevel], flag),
		fatalLevel:   log.New(errOut, f.GetPrefix(fatalLevel)+prefix[fatalLevel], flag),
	}

	return &Logger{formatter: f, loggers: loggers}
}

// Info logs to INFO level
func (l *Logger) Info(v ...interface{}) {
	l.commonPrint(infoLevel, v...)
}

// Infof logs a formatted string to INFO level
func (l *Logger) Infof(format string, v ...interface{}) {
	l.commonPrintf(infoLevel, format, v...)
}

// Warning logs to WARNING level
func (l *Logger) Warning(v ...interface{}) {
	l.commonPrint(warningLevel, v...)
}

// Warningf logs a formatted string to WARNING level
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.commonPrintf(warningLevel, format, v...)
}

// Error logs to ERROR level
func (l *Logger) Error(v ...interface{}) {
	l.commonPrint(errorLevel, v...)
}

// Errorf logs a formatted string to ERROR level
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.commonPrintf(errorLevel, format, v...)
}

// Fatal logs to FATAL level
func (l *Logger) Fatal(v ...interface{}) {
	l.commonPrint(fatalLevel, v...)
}

// Fatalf logs a formatted string to FATAL level
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.commonPrintf(fatalLevel, format, v...)
}

func (l *Logger) commonPrint(lvl level, v ...interface{}) {
	v = l.formatter.Format(lvl, v...)
	v = append(v, l.formatter.GetSuffix(lvl))
	l.loggers[lvl].Print(v...)
}

func (l *Logger) commonPrintf(lvl level, format string, v ...interface{}) {
	suffix := l.formatter.GetSuffix(lvl)
	v = l.formatter.Format(lvl, v...)
	l.loggers[lvl].Printf("%s"+format+suffix, v...)
}
