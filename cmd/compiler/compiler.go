package compiler

import (
	"fmt"
	"os"
)

// the main function for the compiler, the entry point to the program
func main() {
	args := os.Args[1:]

	fmt.Println(args)

}
