// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestInitRootCmd(t *testing.T) {
	InitRootCmd()
	assert.Equal(t, "localhost", hostname)
	assert.Equal(t, 5000, port)
}
