package main

import (
	"jack-jack-go/lib/vmtranslator"
	"os"
)

// the main function for the compiler, the entry point to the program
func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		workingDirectory, err := os.Getwd()
		if err == nil {
			vmtranslator.Translate(workingDirectory)
		}
	} else {
		for _, arg := range args {
			vmtranslator.Translate(arg)
		}
	}
}
