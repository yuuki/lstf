//go:generate go-bindata -nometadata -pkg main -o credits.go CREDITS
package main

import (
	"encoding/json"
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

var (
	creditsText = string(MustAsset("CREDITS"))
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
		numeric bool
		json    bool
		ver     bool
		credits bool
	)
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.BoolVar(&numeric, "n", false, "")
	flags.BoolVar(&numeric, "numeric", false, "")
	flags.BoolVar(&json, "json", false, "")
	flags.BoolVar(&ver, "version", false, "")
	flags.BoolVar(&credits, "credits", false, "")
	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeErr
	}

	if ver {
		fmt.Fprintf(c.errStream, "%s version %s, build %s, date %s \n", name, version, commit, date)
		return exitCodeOK
	}

	if credits {
		fmt.Fprintln(c.outStream, creditsText)
		return exitCodeOK
	}

	flows, err := tcpflow.GetHostFlows()
	if err != nil {
		log.Printf("failed to get host flows: %v", err)
		return exitCodeErr
	}

	if json {
		if err := c.PrintHostFlowsAsJSON(flows, numeric); err != nil {
			log.Printf("failed to print json: %v", err)
			return exitCodeErr
		}
	} else {
		c.PrintHostFlows(flows, numeric)
	}

	return exitCodeOK
}

// PrintHostFlows prints the host flows.
func (c *CLI) PrintHostFlows(flows tcpflow.HostFlows, numeric bool) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(c.outStream, 0, 8, 0, '\t', 0)
	fmt.Fprintln(tw, "Local Address:Port\t <--> \tPeer Address:Port \tConnections")
	for _, flow := range flows {
		if !numeric {
			flow.ReplaceLookupedName()
		}
		fmt.Fprintf(tw, "%s\n", flow)
	}
	tw.Flush()
}

// PrintHostFlowsAsJSON prints the host flows as json format.
func (c *CLI) PrintHostFlowsAsJSON(flows tcpflow.HostFlows, numeric bool) error {
	for _, flow := range flows {
		if !numeric {
			flow.ReplaceLookupedName()
		}
	}
	b, err := json.Marshal(flows)
	if err != nil {
		return err
	}
	c.outStream.Write(b)
	return nil
}

var helpText = `Usage: lstf [options]

  Print host flows between localhost and other hosts

Options:
  --numeric, -n             show numerical addresses instead of trying to determine symbolic host names.
  --json                    print results as json format
  --version, -v	            print version
  --help, -h                print help
  --credits                 print CREDITS
`
