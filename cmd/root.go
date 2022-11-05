// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	defaultHostname = "localhost" // "data.salad.com"
	defaultPort     = 5000
	waitTime        = 15
)

var (
	help     = flag.Bool("h", false, "Show command usage")
	verbose  = flag.Bool("v", false, "Show all the bits")
	hostname string
	port     int
)

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

func Execute() {
	// Resolve hostname and port
	serverEndpoint := fmt.Sprintf("%s:%d", hostname, port)
	serverAddr, err := net.ResolveTCPAddr("tcp", serverEndpoint)
	if err != nil {
		fmt.Printf("Hostname not found: %s\n", hostname)
		fmt.Println(err)
		os.Exit(2)
	}

	log.Printf("Endpoint: %s:%d (%s)\n", hostname, port, serverAddr)

	for {
		// Connect
		log.Println("Connecting to ", serverAddr)
		conn, err := net.DialTCP("tcp", nil, serverAddr)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				// Hang out a bit and try again
				log.Println("Connection refused, pausing for retry")
				time.Sleep(waitTime * time.Second)
				continue
			}
			fmt.Printf("Connect failed:\n")
			fmt.Println(err)
			os.Exit(2)
		}

		// buffer to get data
		log.Println("Reading data")
		var receiveBuf bytes.Buffer
		io.Copy(&receiveBuf, conn)
		fmt.Printf("buflen: %d\n", receiveBuf.Len())
		println("Received message:", receiveBuf.String())

		conn.Close()

		// decode

		// We're retrying too fast, pause a bit...
		time.Sleep(waitTime * time.Second)
	}
}
