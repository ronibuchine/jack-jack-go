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
// expects node of type parameter list
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

	// constant terms
	firstChild := node.Children[0]
	if firstChild.Token.Kind == fe.KEYWORD {
		switch keyword := firstChild.Token.Contents; keyword {
		case "true":
			j.vmw.WritePush("constant", "1")
			j.vmw.w.WriteString("neg\n")
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
		j.compileTerm(node.Children[1])
		switch firstChild.Token.Contents {
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
			switch symbol.kind {
			case "field":
				j.vmw.WritePush("this", strconv.Itoa(symbol.id))
			case "arg":
				j.vmw.WritePush("argument", strconv.Itoa(symbol.id))
			default:
				j.vmw.WritePush(symbol.kind, strconv.Itoa(symbol.id))
			}
			if len(node.Children) > 1 {
				switch node.Children[1].Token.Contents {
				case ".": // method
					j.vmw.WriteCall(symbol.vType+"."+node.Children[2].Token.Contents, j.compileExpressionList(node.Children[4])+1) // +1 to account for the object reference being passed
				case "[": // array
					j.compileExpression(node.Children[2])
					j.vmw.WriteArithmetic("+")
					j.vmw.WritePop("pointer", 1)
					j.vmw.WritePush("that", "0")
				}
			}
		} else { // subroutineCall
			if node.Children[1].Token.Contents == "." { // function or constructor
				j.vmw.WriteCall(firstChild.Token.Contents+"."+node.Children[2].Token.Contents, j.compileExpressionList(node.Children[4]))
			} else {
				// if method, push "this" to stack
				if _, err := j.localST.Find("this"); err != nil {
					j.vmw.WritePush("argument", "0")
					j.vmw.WriteCall(j.className+"."+firstChild.Token.Contents, j.compileExpressionList(node.Children[2])+1)
				} else { // function
					j.vmw.WriteCall(j.className+"."+firstChild.Token.Contents, j.compileExpressionList(node.Children[2]))
				}
			}
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
func (j *JackCompiler) compileLet(node *fe.Node) {
	if symbol, err := j.findSymbol(node.Children[1].Token.Contents); err == nil {
		for count, eq := range node.Children {
			if eq.Token.Contents == "=" {
				j.compileExpression(node.Children[count+1])
			}
		}
		switch symbol.kind {
		case "field":
			j.vmw.WritePop("this", symbol.id)
		case "arg":
			j.vmw.WritePop("argument", symbol.id)
		default:
			j.vmw.WritePop(symbol.kind, symbol.id)
		}
	}
}

// this can be just compileExpression and then pop the return value away
// expects node of kind doStatement
func (j *JackCompiler) compileDo(node *fe.Node) {
	j.compileExpression(&fe.Node{Children: node.Children[1:]})
	j.vmw.WritePop("temp", 0)
}

// expects node of kind ifStatement
// var ifCounter = 0

func (j *JackCompiler) compileIf(node *fe.Node) {
	/* counter := ifCounter
	ifCounter++ */

    endLabel := j.vmw.NewLabel("ifEnd")
    trueLabel := j.vmw.NewLabel("ifTrue")

	j.compileExpression(node.Children[2])
	j.vmw.WriteGoto(trueLabel)
	j.compileStatements(node.Children[9])
	j.vmw.WriteGoto(endLabel)
	j.vmw.WriteLabel(trueLabel)
	j.compileStatements(node.Children[5])
	j.vmw.WriteLabel(endLabel)
}

// expects node of kind whileStatement
var whileCounter = 0

func (j *JackCompiler) compileWhile(node *fe.Node) {
	counter := whileCounter
	whileCounter++
	j.vmw.WriteLabel(j.className + "WhileBegin" + strconv.Itoa(counter))
	j.compileExpression(node.Children[2])
	j.vmw.WriteArithmetic("not")
	j.vmw.WriteGoto(j.className + "WhileEnd" + strconv.Itoa(counter))
	j.compileStatements(node.Children[5])
	j.vmw.WriteGoto(j.className + "WhileBegin" + strconv.Itoa(counter))
	j.vmw.WriteLabel(j.className + "WhileEnd" + strconv.Itoa(counter))
}

// expects node of kind return statement
func (j *JackCompiler) compileReturn(node *fe.Node) {
	if len(node.Children) > 2 {
		j.compileExpression(node.Children[1])
	}
	j.vmw.WriteReturn()
}
