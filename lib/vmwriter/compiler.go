package vmwriter

import (
	"bufio"
	fe "jack-jack-go/lib/syntaxAnalyzer"
	"log"
)

type JackCompiler struct {
	st   *SymbolTable
	ast  *fe.Node
	vmw  *VMWriter
	name string
}

func NewJackCompiler(ast *fe.Node, name string, w *bufio.Writer) *JackCompiler {
	return &JackCompiler{
		ast:  ast,
		vmw:  NewVMWriter(name, w),
		name: name,
	}
}

func (j *JackCompiler) Compile() {
    var varDecs, subRoutines []*fe.Node
    className := j.ast.Children[1];
    if className.Token.Kind != j.name {
        log.Fatal("class name must match the file name")
    }
    
    for _, n := range j.ast.Children {
        if n.Token.Kind == "classVarDec" {
            varDecs = append(varDecs, n)
        } else if n.Token.Kind == "subroutineDec" {
            subRoutines = append(subRoutines, n)
        }
    }

    // assign class level symbol table
    j.st = ClassTable(varDecs)

    for _, n := range subRoutines {
        j.compileSubroutine(n)
    }
}

// expressions
func (j *JackCompiler) compileExpression(node *fe.Node) {
}

func (j *JackCompiler) compileString(node *fe.Node) {
}

// arr[expr1] = expr2?
func (j *JackCompiler) compileArray(node *fe.Node) {
}

// functions
func (j *JackCompiler) compileSubroutine(node *fe.Node) {
    switch node.Token.Kind {
    case "constructor":

    case "function":
    case "method":
    }
}

// returns the local symbol table
func (j *JackCompiler) compileParameterList(node *fe.Node) *SymbolTable {
    return nil
}

func (j *JackCompiler) compileSubroutineBody(node *fe.Node) {
    // first add all local variables to the local symbol table
    // then compile all statements in the subroutine
}

func (j *JackCompiler) compileVarDec(node *fe.Node) {
}

func (j *JackCompiler) compileTerm(node *fe.Node) {
}

func (j *JackCompiler) compileExpressionList(node *fe.Node) {
}

// statements
func (j *JackCompiler) compileLet(node *fe.Node) {
}

// this can be just compileExpression and then pop the return value away
func (j *JackCompiler) compileDo(node *fe.Node) {
}

func (j *JackCompiler) compileIf(node *fe.Node) {
}

func (j *JackCompiler) compileWhile(node *fe.Node) {
}

func (j *JackCompiler) compileReturn(node *fe.Node) {
}
