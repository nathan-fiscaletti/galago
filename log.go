package flamingo

import (
	"log"
	"os"
)

var printAccess bool = false
var logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags)

// SetLogger sets the logger to use for internal flamingo log messages
func SetLogger(l *log.Logger) {
	logger = l
}

// SetPrintAccess tells the system whether or not it should print
// ACCESS log messages to the logger.
func SetPrintAccess(val bool) {
	printAccess = val
}
