package scenariolib

import (
	"io"
	"log"
)

var (
	// Trace Trace logging level
	Trace *log.Logger

	// Info Info logging level
	Info *log.Logger

	// Warning Warning logging level
	Warning *log.Logger

	// Error Error logging level
	Error *log.Logger
)

// InitLogger Initialize the logger with different io.Writer for the the different logging levels
func InitLogger(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

	Trace = log.New(traceHandle, "TRACE >>> ", log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle, "INFO >>> ", log.Ldate|log.Ltime)

	Warning = log.New(warningHandle, "WARNING >>> ", log.Ldate|log.Ltime)

	Error = log.New(errorHandle, "ERROR >>> ", log.Ldate|log.Ltime|log.Lshortfile)
}
