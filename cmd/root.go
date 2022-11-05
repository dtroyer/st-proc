// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	help    = flag.Bool("h", false, "Show command usage")
	verbose = flag.Bool("v", false, "Show all the bits")
)

func Execute() {
	var err error
	log.SetFlags(0)

	flag.Parse()
	args := flag.Args()

	if *help {
		fmt.Printf("Usage: %s [-v] [<hostname> [<port>]]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var host string
	host = "data.salad.com"
	if len(args) >= 1 {
		host = args[0]
	}

	var port int
	port = 5000
	if len(args) >= 2 {
		// port, err = strconv.ParseInt(args[1], 10, 16)
		_, err = fmt.Sscan(args[1], &port)
		if err != nil {
			// Conversion error for arg
			fmt.Printf("Error parsing port: %s\n", args[1])
			fmt.Println(err)
			os.Exit(1)
		}
	}

	log.Printf("Host: %s:%d\n", host, port)
}
