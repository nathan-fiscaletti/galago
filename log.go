package galago

import (
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags)

// SetLogger sets the logger to use for internal galago log messages
func SetLogger(l *log.Logger) {
	logger = l
}
