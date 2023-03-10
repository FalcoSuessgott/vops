package main

import (
	"fmt"
	"os"

	"github.com/FalcoSuessgott/vops/cmd"
)

var version = ""

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "%v.\n", err)
		os.Exit(1)
	}
}
