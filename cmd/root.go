// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

// Entry points for the processor
// * handle command-line arguments
// * main processing loop

package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dtroyer/st-proc/flight"
)

const (
	defaultHostname = "data.salad.com"
	defaultPort     = 5000
	waitTime        = 3
)

var (
	help     = flag.Bool("h", false, "Show command usage")
	verbose  = flag.Bool("v", false, "Show all the bits")
	hostname string
	port     int
)

// Do command line argument processing and other initialization
func InitRootCmd() {
	var err error
	log.SetFlags(0)

	flag.Parse()
	args := flag.Args()

	if *help {
		fmt.Printf("Usage: %s [-v] [<hostname> [<port>]]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	hostname = defaultHostname
	if len(args) >= 1 {
		hostname = args[0]
	}

	port = defaultPort
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
}

// Main processing loop
// The normal exit for the loop is with SIGINT, SIGTERM or Ctrl-C. The
// RouterConn.Connect() method will continuously retry following connection
// refused (ECONNREFUSED) errors, other errors will exit.
func Execute() {
	// Set up the RouterConn host info and validate name resolution
	router := &RouterConn{Hostname: hostname, Port: port, Wait: waitTime}
	err := router.Setup()
	if err != nil {
		fmt.Printf("Hostname not found: %s\n", hostname)
		fmt.Println(err)
		os.Exit(2)
	}

	for {
		// Connect to router
		log.Println("Connecting to ", hostname, ":", port)
		err := router.Connect()
		if err != nil {
			fmt.Printf("Connect failed:\n")
			fmt.Println(err)
			os.Exit(2)
		}

		// Read message data
		log.Println("Reading data")
		buf, err := router.Read()
		log.Printf(" bytes read: %d\n", buf.Len())

		router.Close()

		// Decode the message
		var flightMsg flight.FlightMessage
		err = flight.DecodePacketBuffer(&buf, &flightMsg)
		if err != nil {
			fmt.Println("error: ", err)
		}

		// Make it easy on the eyes: fancy indented JSON
		jsonMsg, err := json.MarshalIndent(flightMsg, "", "  ")
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Printf("%s\n", jsonMsg)
	}
}
