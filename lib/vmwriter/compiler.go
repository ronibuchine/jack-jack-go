package vmwriter

import (
	"bufio"
	"errors"
	fe "jack-jack-go/lib/syntaxAnalyzer"
	"strconv"
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

func (j *JackCompiler) findSymbol(name string) (symbol TableEntry, err error) {
	symbol, err = j.localST.Find(name)
	if err == nil {
		return symbol, nil
	}
	return j.classST.Find(name)
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
		case "static":
			err = j.classST.Add("static", vType, name)
		case "field":
			err = j.classST.Add("this", vType, name)
		case "var":
			err = j.classST.Add("local", vType, name)
		case "arg":
			err = j.localST.Add("argument", vType, name)
		}
		if err != nil {
			return formatError(node, err)
		}
	}
	return nil
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
	err = j.compileSubroutineBody(node.Children[6], node.Children[2].Token.Contents, node.Children[0].Token.Contents == "constructor")
	if err != nil {
		return err
	}
	return nil
}

// adds all parameters in parameter list to the local symbol table
// expects node of type parameter list
func (j *JackCompiler) compileParameterList(params *fe.Node) error {
	var name, vType string
	for i := 0; i < len(params.Children); i += 3 {
		vType = params.Children[i].Token.Contents
		name = params.Children[i+1].Token.Contents
		err := j.localST.Add("argument", vType, name)
		if err != nil {
			return formatError(params, err)
		}
	}
	return nil
}

// expects node of kind subroutineBody
func (j *JackCompiler) compileSubroutineBody(node *fe.Node, funcName string, constructor bool) error {
	// first add all local variables to the local symbol table
	// then compile all statements in the subroutine
	for _, element := range node.Children[1:] {
		switch element.Token.Kind {
		case "varDec":
			j.compileVarDec(element)
		case "statements":
			j.vmw.WriteFunction(j.className+"."+funcName, j.localST.counts["local"])
			if constructor {
				j.vmw.WritePush("constant", strconv.Itoa(j.classST.counts["this"]))
				j.vmw.WriteCall("Memory.alloc", 1)
				j.vmw.WritePop("pointer", "0")
			}
			if _, err := j.localST.Find("this"); err == nil { // if method
				j.vmw.WritePush("argument", "0")
				j.vmw.WritePop("pointer", "0")
			}
			j.compileStatements(element)
		}

	}
	return nil
}

// expects node of kind statements
func (j *JackCompiler) compileStatements(node *fe.Node) {
	for _, statement := range node.Children {
		switch statement.Token.Kind {
		case "ifStatement":
			j.compileIf(statement)
		case "whileStatement":
			j.compileWhile(statement)
		case "doStatement":
			j.compileDo(statement)
		case "letStatement":
			j.compileLet(statement)
		case "returnStatement":
			j.compileReturn(statement)
		}
	}
}

// expects node of kind expression
func (j *JackCompiler) compileExpressionList(node *fe.Node) int {
	for _, expression := range node.Children {
		j.compileExpression(expression)
	}
	return len(node.Children)
}

// expects node of kind letStatement
func (j *JackCompiler) compileLet(node *fe.Node) error {
	symbol, err := j.findSymbol(node.Children[1].Token.Contents)
	if err != nil {
		return err
	}
	j.compileExpression(node.Children[3])
	if node.Children[2].Token.Contents == "[" {
		j.vmw.WritePush(symbol.kind, strconv.Itoa(symbol.id))
		j.vmw.WriteArithmetic("+")
		j.compileExpression(node.Children[6])
		j.vmw.WritePop("temp", "0")
		j.vmw.WritePop("pointer", "1")
		j.vmw.WritePush("temp", "0")
		j.vmw.WritePop("that", "0")
	} else {
		j.vmw.WritePop(symbol.kind, strconv.Itoa(symbol.id))
	}
	return nil
}

// this can be just compileExpression and then pop the return value away
// expects node of kind doStatement
func (j *JackCompiler) compileDo(node *fe.Node) {
	j.compileTerm(&fe.Node{Children: node.Children[1 : len(node.Children)-1]})
	j.vmw.WritePop("temp", "0")
}

func (j *JackCompiler) compileIf(node *fe.Node) {

	endLabel := j.vmw.NewLabel("ifEnd")
	trueLabel := j.vmw.NewLabel("ifTrue")

	j.compileExpression(node.Children[2])
	j.vmw.WriteIf(trueLabel)
	j.compileStatements(node.Children[9])
	j.vmw.WriteGoto(endLabel)
	j.vmw.WriteLabel(trueLabel)
	j.compileStatements(node.Children[5])
	j.vmw.WriteLabel(endLabel)
}

func (j *JackCompiler) compileWhile(node *fe.Node) {
	beginLabel := j.vmw.NewLabel("whileBegin")
	endLabel := j.vmw.NewLabel("whileEnd")

	j.vmw.WriteLabel(beginLabel)
	j.compileExpression(node.Children[2])
	j.vmw.WriteArithmetic("not")
	j.vmw.WriteIf(endLabel)
	j.compileStatements(node.Children[5])
	j.vmw.WriteGoto(beginLabel)
	j.vmw.WriteLabel(endLabel)
}

// expects node of kind returnStatement
func (j *JackCompiler) compileReturn(node *fe.Node) {
	if len(node.Children) == 2 {
		j.vmw.WritePush("constant", "0")
	} else {
		j.compileExpression(node.Children[1])
	}
	j.vmw.WriteReturn()
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
	j.vmw.WritePush("constant", strconv.Itoa(len(node.Token.Contents)))
	j.vmw.WriteCall("String.new", 1)
	for _, char := range node.Token.Contents {
		j.vmw.WritePush("constant", strconv.Itoa(int(char)))
		j.vmw.WriteCall("String.appendChar", 2)
	}
}

// expects node of kind term
func (j *JackCompiler) compileTerm(node *fe.Node) {

	// constant terms
	firstChild := node.Children[0]
	if firstChild.Token.Kind == fe.KEYWORD {
		switch keyword := firstChild.Token.Contents; keyword {
		case "true":
			j.vmw.WritePush("constant", "1")
			j.vmw.WriteArithmetic("neg")
		case "false", "null":
			j.vmw.WritePush("constant", "0")
		case "this":
			j.vmw.WritePush("pointer", "0")
		}
	}
	if firstChild.Token.Kind == fe.INT {
		j.vmw.WritePush("constant", firstChild.Token.Contents)
	}
	if firstChild.Token.Kind == fe.STRING {
		j.compileString(firstChild)
	}

	// unary operator
	if firstChild.Token.Kind == fe.SYMBOL {
		j.compileExpression(node.Children[1])
		switch firstChild.Token.Contents {
		// implicitly handles (expression)
		case "-":
			j.vmw.WriteArithmetic("neg")
		case "~":
			j.vmw.WriteArithmetic("not")
		}
	}

	// identifiers
	if firstChild.Token.Kind == fe.IDENT {
		// variable
		if symbol, err := j.findSymbol(firstChild.Token.Contents); err == nil {
			j.vmw.WritePush(symbol.kind, strconv.Itoa(symbol.id))
			if len(node.Children) > 1 {
				switch node.Children[1].Token.Contents {
				case ".": // method
					j.vmw.WriteCall(symbol.vType+"."+node.Children[2].Token.Contents, j.compileExpressionList(node.Children[4])+1) // +1 to account for the object reference being passed
				case "[": // array element
					j.compileExpression(node.Children[2])
					j.vmw.WriteArithmetic("+")
					j.vmw.WritePop("pointer", "0")
					j.vmw.WritePush("that", "0")
				}
			}
		} else { // subroutine call
			if node.Children[1].Token.Contents == "." { // function or constructor
				j.vmw.WriteCall(firstChild.Token.Contents+"."+node.Children[2].Token.Contents, j.compileExpressionList(node.Children[4]))
			} else {
				// if method, push "this" to stack
				if _, err := j.localST.Find("this"); err != nil {
					j.vmw.WritePush("argument", "0")
					j.vmw.WriteCall(j.className+"."+firstChild.Token.Contents, j.compileExpressionList(node.Children[2])+1)
				} else { // function calling function
					j.vmw.WriteCall(j.className+"."+firstChild.Token.Contents, j.compileExpressionList(node.Children[2]))
				}
			}
		}
	}
}
