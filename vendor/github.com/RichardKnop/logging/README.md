[![Codeship Status for AreaHQ/logging](https://codeship.com/projects/a2b1c4c0-a3a2-0133-00d6-3641d785a31d/status?branch=master)](https://codeship.com/projects/129304)

# Logging

A simple leveled logging library with coloured output.

Log levels:
* `INFO` (blue)
* `WARNING` (pink)
* `ERROR` (red)
* `FATAL` (red)

Formatters:
* `DefaultFormatter`
* `ColouredFormatter`

Example usage:

```go
package main

import (
	"github.com/AreaHQ/logging"
)

var (
	plainLogger    *logging.Logger
	colouredLogger *logging.Logger
)

func init() {
	plainLogger = logging.New(nil, nil, nil)
	colouredLogger = logging.New(nil, nil, new(logging.ColouredFormatter))
}

func main() {
	plainLogger.Info("log message")
	plainLogger.Infof("formatted %s %s", "log", "message")
	colouredLogger.Info("log message")
	colouredLogger.Infof("formatted %s %s", "log", "message")

	plainLogger.Warning("log message")
	plainLogger.Warningf("formatted %s %s", "log", "message")
	colouredLogger.Warning("log message")
	colouredLogger.Warningf("formatted %s %s", "log", "message")

	plainLogger.Error("log message")
	plainLogger.Errorf("formatted %s %s", "log", "message")
	colouredLogger.Error("log message")
	colouredLogger.Errorf("formatted %s %s", "log", "message")

	// Not that logger.Fatal/f does not exit program execution
	// To emulate log.Fatal/f, follow with os.Exit(1)
	plainLogger.Fatal("log message")
	plainLogger.Fatalf("formatted %s %s", "log", "message")
	colouredLogger.Fatal("log message")
	colouredLogger.Fatalf("formatted %s %s", "log", "message")
}
```
