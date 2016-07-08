package logging

import (
	"fmt"
)

const (
	// For colouring
	resetSeq  = "\033[0m"
	colourSeq = "\033[0;%dm"
)

// Colour map
var colour = []string{
	infoLevel:    fmt.Sprintf(colourSeq, 94), // blue
	warningLevel: fmt.Sprintf(colourSeq, 95), // pink
	errorLevel:   fmt.Sprintf(colourSeq, 91), // red
	fatalLevel:   fmt.Sprintf(colourSeq, 91), // red
}

// ColouredFormatter colours log messages with ASCI escape codes
// and adds filename and line number before the log message
// See https://en.wikipedia.org/wiki/ANSI_escape_code
type ColouredFormatter struct {
}

// GetPrefix returns colour escape code
func (f *ColouredFormatter) GetPrefix(lvl level) string {
	return colour[lvl]
}

// GetSuffix returns reset sequence code
func (f *ColouredFormatter) GetSuffix(lvl level) string {
	return resetSeq
}

// Format adds filename and line number before the log message
func (f *ColouredFormatter) Format(lvl level, v ...interface{}) []interface{} {
	return append([]interface{}{header()}, v...)
}
