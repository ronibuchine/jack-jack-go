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

func (j *JackCompiler) Compile() error {
	var varDecs, subRoutines []*fe.Node
	className := j.ast.Children[1]
	if className.Token.Contents != j.name {
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
    classST, err := ClassTable(varDecs)
    if err != nil {
        return err
    }
    j.st = classST

	for _, n := range subRoutines {
		j.compileSubroutine(n)
	}

    return nil
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
func (j *JackCompiler) compileSubroutine(node *fe.Node) error {
    st, err := LocalTable(node)
    if err != nil {
        return err
    }
	switch node.Token.Kind {
	case "constructor":
	case "function":
	case "method":
	}
    return nil
}

// returns the local symbol table
// ahhh wtf!!
// I need to just create a parameter list, but if it is a ctor or method then add `this`.
// Where
func (j *JackCompiler) compileParameterList(node *fe.Node, functionKind string) *SymbolTable {
	lst := newSymbolTable()
	var name, vType string
	if functionKind := subroutine.Children[0].Token.Kind; functionKind == "constructor" || functionKind == "method" {
		vType = subroutine.Children[1].Token.Contents
		if vType == "void" && functionKind == "constructor" {
			return nil, errors.New("void constructor makes no sense")
		}
		lst.Add("this", "arg", vType)
	}
	var params *fe.Node
	for _, child := range subroutine.Children {
		if child.Token.Kind == "parameterList" {
			params = child
			break
		}
	}
	for i := 0; i < len(params.Children); i += 3 {
		vType = params.Children[i].Token.Contents
		name = params.Children[i+1].Token.Contents
		err := lst.Add(name, "arg", vType)
		if err != nil {
			return nil, formatError(subroutine, err)
		}
	}
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
