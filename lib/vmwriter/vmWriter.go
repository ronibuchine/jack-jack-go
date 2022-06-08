package vmwriter

import (
	"bufio"
	"fmt"
)

type UniqueLabel string

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

func (vmw *VMWriter) FlushVMFile() {
    vmw.w.Flush()
}

func (vmw *VMWriter) WritePush(segment string, index string) {
	vmw.w.WriteString(fmt.Sprintf("push %s %s\n", segment, index))
}

func (vmw *VMWriter) WritePop(segment string, index int) {
	vmw.w.WriteString(fmt.Sprintf("pop %s %d\n", segment, index))
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
	case "&":
		vmw.w.WriteString("and\n")
	case "|":
		vmw.w.WriteString("or\n")
	case "<":
		vmw.w.WriteString("lt\n")
	case ">":
		vmw.w.WriteString("gt\n")
	case "eq":
		vmw.w.WriteString("eq\n")
	case "neg":
		vmw.w.WriteString("neg\n")
	case "not":
		vmw.w.WriteString("not\n")
	}
}

func (vmw *VMWriter) NewLabel(label string) UniqueLabel {
	s := fmt.Sprintf("%s_%s_%d", vmw.name, label, vmw.labelCounter)
	vmw.labelCounter++
	return UniqueLabel(s)
}

// returns the unique label that was written
func (vmw *VMWriter) WriteLabel(label UniqueLabel) {
	vmw.w.WriteString(fmt.Sprintf("label %s\n", label))
}

func (vmw *VMWriter) WriteGoto(label UniqueLabel) {
	vmw.w.WriteString(fmt.Sprintf("goto %s\n", label))
}

func (vmw *VMWriter) WriteIf(label UniqueLabel) {
	vmw.w.WriteString(fmt.Sprintf("if-goto %s\n", label))
}

func (vmw *VMWriter) WriteCall(name string, nargs int) {
	vmw.w.WriteString(fmt.Sprintf("call %s %d\n", name, nargs))
}

func (vmw *VMWriter) WriteFunction(name string, nargs int) {
	vmw.w.WriteString(fmt.Sprintf("function %s %d\n", name, nargs))
}

func (vmw *VMWriter) WriteReturn() {
	vmw.w.WriteString("return\n")
}
