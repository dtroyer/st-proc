// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

// RouterConn handles the network I/O to the message router

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"
)

// RouterConn is a simple receive-only connection to the message router
// * infinite retry if connection is refused
// * reads all sent data until server closes connection

type RouterConn struct {
	Hostname   string
	Port       int
	Wait       int
	serverAddr *net.TCPAddr
	connection *net.TCPConn
}

func (m *RouterConn) Close() {
	m.connection.Close()
}

// Connect to router, retry on ECONNREFUSED, abort on all other errors
func (m *RouterConn) Connect() (err error) {
	for {
		m.connection, err = net.DialTCP("tcp", nil, m.serverAddr)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				// Hang out a bit and try again
				time.Sleep(time.Duration(m.Wait) * time.Second)
				continue
			}
			return err
		}
		break
	}
	return nil
}

// Read from network until server closes connection
// TODO: what is the max size???
func (m *RouterConn) Read() (receiveBuf bytes.Buffer, err error) {
	io.Copy(&receiveBuf, m.connection)
	return receiveBuf, nil
}

// Set up the address struct and do hostname resolution
func (m *RouterConn) Setup() (err error) {
	m.serverAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", m.Hostname, m.Port))
	return err
}
