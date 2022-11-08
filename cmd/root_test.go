// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

// Test command-line defaults
func TestInitRootCmd(t *testing.T) {
	InitRootCmd()
	assert.Equal(t, defaultHostname, hostname)
	assert.Equal(t, defaultPort, port)
}
