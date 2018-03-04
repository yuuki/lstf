package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
