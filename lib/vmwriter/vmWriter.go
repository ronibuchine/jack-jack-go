package vmwriter

import (
	"bufio"
)

type VMWriter struct {
    w *bufio.Writer
    labelCounter int
    name string
}

func NewVMWriter(name string, w *bufio.Writer) *VMWriter {
    return nil
}

func (vmw *VMWriter) WritePush(segment string, index int) {
}

func (vmw *VMWriter) WritePop(segment string, index int) {
}

func (vmw *VMWriter) WriteArithmetic(command string) {
}

// returns the unique label that was written
func (vmw *VMWriter) WriteLabel() string {
    return ""
}

func (vmw *VMWriter) WriteGoto() {
}

func (vmw *VMWriter) WriteIf() {
}

func (vmw *VMWriter) WriteCall(name string, nargs int) {
}

func (vmw *VMWriter) WriteFunction(name string, nargs int) {
}

func (vmw *VMWriter) WriteReturn() {
}
