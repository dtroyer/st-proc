// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/dtroyer/st-proc/cmd"
)

func main() {
	cmd.InitRootCmd()
	cmd.Execute()
}
