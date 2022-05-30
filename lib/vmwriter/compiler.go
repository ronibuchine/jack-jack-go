package vmwriter

import (
    fe "jack-jack-go/lib/syntaxAnalyzer"
)

// any function that may contain an expression needs to see the class level symbol table

// classes dont actually contain any code. They are just collections of
// variables and functions
func CompileClass(node *fe.Node) {
}

// this doesn't write any code, just returns a symbol table
func CompileClassVarDec(node *fe.Node) *SymbolTable {
    return nil
}

// expressions
func CompileExpression(node *fe.Node, className string, st *SymbolTable) {
}

func CompileString(node *fe.Node) {
}

// arr[expr1] = expr2?
func CompileArray(node *fe.Node, st *SymbolTable) {
}


// functions
func CompileSubroutine(node *fe.Node, className string, st *SymbolTable) {
}

func CompileParameterList(node *fe.Node) {
}

func CompileSubroutineBody(node *fe.Node, st *SymbolTable) {
}

func CompileVarDec(node *fe.Node) {
}

func CompileTerm(node *fe.Node, st *SymbolTable) {
}

func CompileExpressionList(node *fe.Node, st *SymbolTable) {
}

// statements
func CompileLet(node *fe.Node, st *SymbolTable) {
}

// this can be just compileExpression and then pop the return value away
func CompileDo(node *fe.Node, st *SymbolTable) {
}


func CompileIf(node *fe.Node, st *SymbolTable) {
}

func CompileWhile(node *fe.Node, st *SymbolTable) {
}

func CompileReturn(node *fe.Node, st *SymbolTable) {
}

