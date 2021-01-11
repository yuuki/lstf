//go:generate go-bindata -nometadata -pkg main -o credits.go CREDITS
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/yuuki/lstf/dlog"
	"github.com/yuuki/lstf/tcpflow"
	"golang.org/x/xerrors"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 10 + iota
)

var (
	creditsText = string(MustAsset("CREDITS"))
)

func setDebugOutputLevel(debug bool) {
	if debug {
		dlog.Debug = true
	}

	debugEnv := os.Getenv("LSTF_DEBUG")
	if debugEnv != "" {
		showDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing boolean value from LSTF_DEBUG: %s\n", err)
			os.Exit(1)
		}
		dlog.Debug = showDebug
	}
}

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
		numeric   bool
		processes bool
		watch     time.Duration
		json      bool
		filter    string

		ver     bool
		credits bool
		debug   bool
	)
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.BoolVar(&numeric, "n", false, "")
	flags.BoolVar(&numeric, "numeric", false, "")
	flags.BoolVar(&processes, "p", false, "")
	flags.BoolVar(&processes, "processes", false, "")
	flags.DurationVar(&watch, "w", -1, "")
	flags.DurationVar(&watch, "watch", -1, "")
	flags.BoolVar(&json, "json", false, "")
	flags.StringVar(&filter, "f", tcpflow.FilterAll, "")
	flags.StringVar(&filter, "filter", tcpflow.FilterAll, "")
	flags.BoolVar(&ver, "version", false, "")
	flags.BoolVar(&credits, "credits", false, "")
	flags.BoolVar(&debug, "debug", false, "")
	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeErr
	}

	setDebugOutputLevel(debug)

	if ver {
		fmt.Fprintf(c.errStream, "%s version %s, build %s, date %s \n", name, version, commit, date)
		return exitCodeOK
	}

	if credits {
		fmt.Fprintln(c.outStream, creditsText)
		return exitCodeOK
	}

	if !(filter == tcpflow.FilterAll ||
		filter == tcpflow.FilterPublic ||
		filter == tcpflow.FilterPrivate) {
		fmt.Fprint(c.errStream, helpText)
		return exitCodeErr
	}

	if watch == -1 {
		return c.run(processes, numeric, json, filter)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	stopChan := make(chan struct{})
	retChan := make(chan int, 1)

	go func() {
		tick := time.NewTicker(watch)
		for {
			select {
			case <-tick.C:
				ret := c.run(processes, numeric, json, filter)
				if ret != exitCodeOK {
					retChan <- ret
					break
				}
			case <-stopChan:
				retChan <- exitCodeOK
				break
			default:
			}
		}
	}()

	<-sig
	stopChan <- struct{}{}
	return <-retChan
}

func (c *CLI) run(processes, numeric, json bool, filter string) int {
	flows, err := tcpflow.GetHostFlows(&tcpflow.GetHostFlowsOption{
		Processes: processes,
		Filter:    filter,
		Numeric:   numeric,
	})
	if err != nil {
		if dlog.Debug {
			log.Printf("failed to get host flows: %+v\n", err)
		} else {
			log.Printf("failed to get host flows: %v\n", err)
		}
		return exitCodeErr
	}

	if json {
		if err := c.PrintHostFlowsAsJSON(flows); err != nil {
			log.Printf("failed to print json: %v\n", err)
			return exitCodeErr
		}
	} else {
		c.PrintHostFlows(flows, processes)
	}

	return exitCodeOK
}

// PrintHostFlows prints the host flows.
func (c *CLI) PrintHostFlows(flows tcpflow.HostFlows, processes bool) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(c.outStream, 0, 8, 0, '\t', 0)
	fmt.Fprintf(tw, "Local Address:Port\t<-->\tPeer Address:Port\tConnections")
	if processes {
		fmt.Fprintf(tw, "\tProcess")
	}
	fmt.Fprintln(tw)
	for _, flow := range flows {
		fmt.Fprintln(tw, flow)
	}
	tw.Flush()
}

// PrintHostFlowsAsJSON prints the host flows as json format.
func (c *CLI) PrintHostFlowsAsJSON(flows tcpflow.HostFlows) error {
	b, err := json.Marshal(flows)
	if err != nil {
		return xerrors.Errorf("failed to marshal json: %v", err)
	}
	c.outStream.Write(b)
	return nil
}

var helpText = `Usage: lstf [options]

  Print host flows between localhost and other hosts

Options:
  --numeric, -n             show numerical addresses instead of trying to determine symbolic host names.
  --processes, -p           show process using socket
  --json                    print results as json format
  --filter, -f FILTER       filter results by "all", "public" or "private" (default: "all")
  --watch, -w DURATION      print periodically (DURATION should be like '3s')
  --version, -v	            print version
  --help, -h                print help
  --credits                 print CREDITS
`
