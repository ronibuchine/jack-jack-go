package vmwriter

import (
	"bufio"
)

// embedded types ftw
type VMWriter struct {
    bufio.Writer
    code string
    labelCounter int
    name string
}

func WritePush(segment string, index int) {
}

func WritePop(segment string, index int) {
}

func WriteArithmetic(command string) {
}

// returns the unique label that was written
func WriteLabel() string {
    return ""
}

func WriteGoto() {
}

func WriteIf() {
}

func WriteCall(name string, nargs int) {
}

func WriteFunction(name string, nargs int) {
}

func WriteReturn() {
}
