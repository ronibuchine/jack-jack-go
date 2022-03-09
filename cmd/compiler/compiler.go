package compiler

import (
	"cmd/lib/vmtranslator"
	"os"
)

// the main function for the compiler, the entry point to the program
func main() {

	args := os.Args[1:]
	for _, arg := range args {
		vmtranslator.Translate(arg)
	}
}
