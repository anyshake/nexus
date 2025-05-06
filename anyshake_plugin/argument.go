package main

import (
	"flag"
)

type arguments struct {
	address string
	timeout int
	verbose bool
}

func parseCommandLine() (args arguments) {
	flag.StringVar(&args.address, "address", "127.0.0.1:30000", "Address of TCP forwarder service defined in AnyShake Observer")
	flag.IntVar(&args.timeout, "timeout", 10, "Timeout value of TCP forwarder connection in seconds")

	var verbose string
	flag.StringVar(&verbose, "verbose", "false", "Enable verbose logging")
	args.verbose = (verbose == "true")

	flag.Parse()
	return args
}
