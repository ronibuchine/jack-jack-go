package main

import (
	"os"
)

const OUTPUT_FILE_NAME string = "Tar0.asm"

func main() {
	os.Create(OUTPUT_FILE_NAME)

}
