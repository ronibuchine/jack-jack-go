package vmwriter

import (
	"bufio"
	"fmt"
)

type VMWriter struct {
	w            *bufio.Writer
	labelCounter int
	name         string
}

func NewVMWriter(name string, w *bufio.Writer) *VMWriter {
	return &VMWriter{
		w:            w,
		labelCounter: 0,
		name:         name,
	}
}

func (vmw *VMWriter) WritePush(segment string, index int) {
	vmw.w.WriteString(fmt.Sprintf("push %s %s\n", segment, index))
}

func (vmw *VMWriter) WritePop(segment string, index int) {
	vmw.w.WriteString(fmt.Sprintf("pop %s %s\n", segment, index))
}

func (vmw *VMWriter) WriteArithmetic(command string) {
	switch command {
	case "+":
		vmw.w.WriteString("add\n")
	case "-":
		vmw.w.WriteString("sub\n")
	case "*":
		vmw.w.WriteString("call Math.multiply 2\n")
	case "/":
		vmw.w.WriteString("call Math.divide 2\n")
	}
}

// returns the unique label that was written
func (vmw *VMWriter) WriteLabel(label string) string {
	vmw.w.WriteString(fmt.Sprintf("label %s\n", label))
	return label
}

func (vmw *VMWriter) WriteGoto(label string) {
	vmw.w.WriteString(fmt.Sprintf("goto %s\n", label))
}

func (vmw *VMWriter) WriteIf(label string) {
	vmw.w.WriteString(fmt.Sprintf("if-goto %s\n", label))
}

func (vmw *VMWriter) WriteCall(name string, nargs int) {
	vmw.w.WriteString(fmt.Sprintf("call %s %s\n", name, nargs))
}

func (vmw *VMWriter) WriteFunction(name string, nargs int) {
	vmw.w.WriteString(fmt.Sprintf("function %s %s\n", name, nargs))
}

func (vmw *VMWriter) WriteReturn() {
	vmw.w.WriteString("return\n")
}
