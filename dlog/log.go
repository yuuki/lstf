package dlog

import (
	"log"
)

// Debug is option variable for debug mode.
var Debug = false

// Debugln prints debug log.
func Debugln(v ...interface{}) {
	if Debug == true {
		origin := log.Flags()
		log.SetFlags(log.LstdFlags | log.Llongfile)
		log.Println(v...)
		log.SetFlags(origin)
	}
}

// Debugf prints debug log with specified format.
func Debugf(format string, v ...interface{}) {
	if Debug == true {
		origin := log.Flags()
		log.SetFlags(log.LstdFlags | log.Llongfile)
		log.Printf(format, v...)
		log.SetFlags(origin)
	}
}
