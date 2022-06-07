package vmwriter

import (
	"bufio"
	"errors"
	fe "jack-jack-go/lib/syntaxAnalyzer"
)

type JackCompiler struct {
	classST   *SymbolTable
	localST   *SymbolTable
	ast       *fe.Node
	vmw       *VMWriter
	className string
}

func NewJackCompiler(ast *fe.Node, name string, w *bufio.Writer) *JackCompiler {
	return &JackCompiler{
		classST:   newSymbolTable(),
		localST:   newSymbolTable(),
		ast:       ast,
		vmw:       NewVMWriter(name, w),
		className: name,
	}
}

func (j *JackCompiler) findSymbolKind(name string) (kind string, err error) {
	kind, err = j.localST.KindOf(name)
	if err == nil {
		return kind, nil
	}
	kind, err = j.classST.KindOf(name)
	if err == nil {
		return kind, nil
	}
	return "", err
}

func (j *JackCompiler) findSymbolType(name string) (vType string, err error) {
	vType, err = j.localST.TypeOf(name)
	if err == nil {
		return vType, nil
	}
	vType, err = j.classST.TypeOf(name)
	if err == nil {
		return vType, nil
	}
	return "", err
}

func (j *JackCompiler) findSymbolIndex(name string) (index int, err error) {
	index, err = j.localST.IndexOf(name)
	if err == nil {
		return index, nil
	}
	index, err = j.classST.IndexOf(name)
	if err == nil {
		return index, nil
	}
	return 0, err
}

func (j *JackCompiler) CompileClass() error {
	className := j.ast.Children[1]
	if className.Token.Contents != j.className {
		return errors.New("class name must match the file name")
	}
	var err error
	for _, n := range j.ast.Children {
		if n.Token.Kind == "classVarDec" {
			err = j.compileVarDec(n)
		} else if n.Token.Kind == "subroutineDec" {
			err = j.compileSubroutine(n)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// node should be of kind expression
func (j *JackCompiler) compileExpression(node *fe.Node) {
	j.compileTerm(node.Children[0])
	for termCount := 1; termCount < len(node.Children); termCount += 2 {
		j.compileTerm(node.Children[termCount+1])
		j.vmw.WriteArithmetic(node.Children[termCount].Token.Contents)
	}
}

// node should be of kind string
func (j *JackCompiler) compileString(node *fe.Node) {
}

// arr[expr1] = expr2?
func (j *JackCompiler) compileArray(node *fe.Node) {
}

// node should be of kind constructor, function or method
func (j *JackCompiler) compileSubroutine(node *fe.Node) (err error) {
	j.localST.Clear()
	if node.Children[0].Token.Contents == "method" {
		j.localST.Add("arg", j.className, "this")
	}
	err = j.compileParameterList(node.Children[4])
	if err != nil {
		return err
	}
	err = j.compileSubroutineBody(node.Children[6])
	if err != nil {
		return err
	}
	if node.Children[0].Token.Contents == "constructor" {
		// add `return this` in vmish
	}
	return nil
}

// adds all parameters in parameter list to the local symbol table
// expects node of of type parameter list
func (j *JackCompiler) compileParameterList(params *fe.Node) error {
	var name, vType string
	for i := 0; i < len(params.Children); i += 3 {
		vType = params.Children[i].Token.Contents
		name = params.Children[i+1].Token.Contents
		err := j.localST.Add("arg", vType, name)
		if err != nil {
			return formatError(params, err)
		}
	}
	return nil
}

// expects node of kind subroutineBody
func (j *JackCompiler) compileSubroutineBody(node *fe.Node) error {
	// first add all local variables to the local symbol table
	// then compile all statements in the subroutine
	return nil
}

// expects node of kind varDec.
// adds it to the appropriate symbol table
func (j *JackCompiler) compileVarDec(node *fe.Node) error {
	var (
		kind, vType, name string
		err               error
	)

	kind = node.Children[0].Token.Contents
	vType = node.Children[1].Token.Contents
	for i := 2; i < len(node.Children); i += 2 {
		name = node.Children[i].Token.Contents
		switch kind {
		case "static", "field":
			err = j.classST.Add(kind, vType, name)
		case "var", "arg":
			err = j.localST.Add(kind, vType, name)
		}
		if err != nil {
			return formatError(node, err)
		}
	}
	return nil
}

// expects node of kind term
func (j *JackCompiler) compileTerm(node *fe.Node) {
}

// expects node of kind expression
func (j *JackCompiler) compileExpressionList(node *fe.Node) {
}

// expects node of kind letStatement
func (j *JackCompiler) compileLet(node *fe.Node) {
}

// this can be just compileExpression and then pop the return value away
// expects node of kind doStatement
func (j *JackCompiler) compileDo(node *fe.Node) {
}

// expects node of kind ifStatement
func (j *JackCompiler) compileIf(node *fe.Node) {
}

// expects node of kind whileStatement
func (j *JackCompiler) compileWhile(node *fe.Node) {
}

// expects node of kind return statement
func (j *JackCompiler) compileReturn(node *fe.Node) {
}
