package vmwriter

import (
    fe "jack-jack-go/lib/syntaxAnalyzer"
)


// classes dont actually contain any code. They are just collections of
// variables (which don't create code) and functions (which do)
func CompileClass(node *fe.Node) {
}

// expressions
func CompileExpression(node *fe.Node, className string) {
}

func CompileString(node *fe.Node) {
}

// arr[expr1] = expr2
func CompileArray(node *fe.Node) {
}


// functions
func CompileConstructor(node *fe.Node, className string) {
}

func CompileMethod(node *fe.Node, className string) {
}

func CompileFunction(node *fe.Node, className string) {
}

// statements
func CompileDo(node *fe.Node) {
}

func CompileLet(node *fe.Node) {
}

func CompileIf(node *fe.Node) {
}

func CompileWhile(node *fe.Node) {
}

func CompileReturn(node *fe.Node) {
}

