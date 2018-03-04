package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"text/tabwriter"

	"github.com/yuuki/lstf/tcpflow"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 10 + iota
)

// CLI is the command line object.
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run execute the main process.
// It returns exit code.
func (c *CLI) Run(args []string) int {
	log.SetOutput(c.errStream)

	var (
		ver bool
	)
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.BoolVar(&ver, "version", false, "")
	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeErr
	}

	if ver {
		fmt.Fprintf(c.errStream, "%s version %s, build %s, date %s \n", name, version, commit, date)
		return exitCodeOK
	}

	flows, err := tcpflow.GetHostFlows()
	if err != nil {
		log.Printf("failed to get host flows: %v", err)
		return exitCodeErr
	}

	c.PrintHostFlows(flows)

	return exitCodeOK
}

// PrintHostFlows prints the host flows.
func (c *CLI) PrintHostFlows(flows tcpflow.HostFlows) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(c.outStream, 0, 8, 0, '\t', 0)
	fmt.Fprintln(tw, "Local Address:Port\t <--> \tPeer Address:Port")
	for _, flow := range flows {
		fmt.Fprintf(tw, "%s\n", flow)
	}
	tw.Flush()
}

var helpText = `Usage: lstf [options]

  Print host flows between localhost and other hosts

Options:
  --version, -v	            print version
  --help, -h                print help
`
