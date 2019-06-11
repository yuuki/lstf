package dlog

import (
	"log"
)

// Debug is option variable for debug mode.
var Debug = false

// Debugln prints debug log.
func Debugln(v ...interface{}) {
	if Debug == true {
		log.Println(v...)
	}
}

// Debugf prints debug log with specified format.
func Debugf(format string, v ...interface{}) {
	if Debug == true {
		log.Printf(format, v...)
	}
}
