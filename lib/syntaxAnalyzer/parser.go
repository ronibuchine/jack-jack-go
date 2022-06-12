package syntaxAnalyzer

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"jack-jack-go/lib/util"
	"log"
	"strings"
)

/*
This file exposes:
    type Node struct -- which is a node in the AST
                     -- Node implements Marshaller which informs xml.Marshal() how to write itself


    func Parse(tokens []Token) *Node -- parses a slice of tokens and returns an AST
    func BuildXML(root *Node, w io.Writer) error -- takes a root node and writes XML to the writer
*/

// Node struct for parsing and recursive descent
type Node struct {
	Token    Token
	Children []*Node
}


func (n Node) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Token.Kind == KEYWORD || n.Token.Kind == SYMBOL ||
		n.Token.Kind == IDENT || n.Token.Kind == INT || n.Token.Kind == STRING {
		return e.Encode(n.Token)
	} else {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Space: "", Local: n.Token.Kind}})
		e.Encode(n.Children)
		return e.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: n.Token.Kind}})
	}
}

func createNodeFromToken(token Token) *Node {
	return &Node{token, []*Node{}}
}

func createNodeFromString(name string) *Node {
	return &Node{Token{Kind: name}, []*Node{}}
}

func (parent *Node) addChild(child *Node) {
	parent.Children = append(parent.Children, child)
}

// parse a stream of tokens and return the root Node of the AST
func Parse(ts *TS) *Node {
	if ts.Tokens[0].Contents != "class" {
		log.Fatal("Jack file must be contained in a class object")
	}
	rootNode := class(ts)
	return rootNode
}

// Build xml and write to disk from root node
func NodeToXML(root *Node, w io.Writer) error {
	bytes, err := xml.MarshalIndent(root, "", "  ")
	// bytes = []byte(xml.Header + string(bytes))
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

var binaryOperators []string = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
var unaryOperators []string = []string{"~", "-"}
var keywordConst []string = []string{"true", "false", "null", "this"}
var functionDecs []string = []string{"function", "constructor", "method"}
var classVars []string = []string{"static", "field"}
var statementKeywords []string = []string{"let", "do", "if", "while", "return"}

func _matchSingle(token string, ts *TS) (*Node, error) {
	// if we match a ident, int or string, we DONT care about the contents
	// if we match a symbol or keyword, we DO care about contents
	curr := ts.curTok()
	if ((token == IDENT || token == INT || token == STRING) && (token == curr.Kind)) ||
		((token == curr.Contents) && (curr.Kind == KEYWORD || curr.Kind == SYMBOL)) {
		res := createNodeFromToken(curr)
        ts.counter++
		return res, nil
	}
	return createNodeFromString("ERROR"), errors.New(fmt.Sprint("Failed to match ", token))
}

// global function used for token parsing and matching
// can either pass in a string or []string. If no matches, then an error will
// occur, if at least one matches then the first match will be returned
func match(token interface{}, ts *TS) (result *Node) {
	if ts.counter >= len(ts.Tokens) {
		log.Fatal("end of token stream")
	}

    switch token := token.(type) {
    case string: 
		if match, err := _matchSingle(token, ts); err == nil {
            return match
		} else {
			parseError(token, ts)
            return createNodeFromString("ERROR")
		}
    case []string:
		for _, t := range token {
			if match, err := _matchSingle(t, ts); err == nil {
                return match
			}
		}
		parseError(strings.Join(token, ", "), ts)
        ts.counter++
        return createNodeFromString("ERROR")
    default:
		panic("match() should only be passed a string or a list of strings")
    }
}

// prints to stdout a parse error
func parseError(expected string, ts *TS) {
	curr := ts.curTok()
    fmt.Print(fmt.Sprintf("ERROR: %s: line %d: Expected token(s) `%s` before %s %s\n",ts.File, curr.LineNumber, expected, curr.Kind, curr.Contents))
}

// functions for grammar
func class(ts *TS) *Node {
	result := createNodeFromString("class")
	result.addChild(match("class", ts))
	result.addChild(match(IDENT, ts))
	result.addChild(match("{", ts))
	curr := ts.curTok()
	for util.Contains(classVars, curr.Contents) {
		result.addChild(classVarDec(ts))
		curr = ts.curTok()
	}
	for util.Contains(functionDecs, curr.Contents) {
		result.addChild(subroutineDec(ts))
		curr = ts.curTok()
	}
	result.addChild(match("}", ts))
	return result
}

func classVarDec(ts *TS) *Node {
	result := createNodeFromString("classVarDec")
	result.addChild(match(classVars, ts))
	result.addChild(typeName(ts))
	result.addChild(match(IDENT, ts))
	for ts.curTok().Contents == "," {
		result.addChild(match(",", ts))
		result.addChild(match(IDENT, ts))
	}
	result.addChild(match(";", ts))
	return result
}

func typeName(ts *TS) *Node {
	result := match([]string{"int", "char", "boolean", IDENT}, ts)
	return result
}

func subroutineDec(ts *TS) *Node {
	result := createNodeFromString("subroutineDec")
	result.addChild(match(functionDecs, ts))
	if ts.curTok().Contents == "void" {
		result.addChild(match("void", ts))
	} else {
		result.addChild(typeName(ts))
	}
	result.addChild(match(IDENT, ts))
	result.addChild(match("(", ts))
	result.addChild(parameterList(ts))
	result.addChild(match(")", ts))
	result.addChild(subroutineBody(ts))
	return result
}

func parameterList(ts *TS) *Node {
	result := createNodeFromString("parameterList")
	if ts.curTok().Contents == ")" {
		return result
	}
	result.addChild(typeName(ts))
	result.addChild(match(IDENT, ts))
	for ts.curTok().Contents == "," {
		result.addChild(match(",", ts))
		result.addChild(typeName(ts))
		result.addChild(match(IDENT, ts))
	}
	return result
}

func subroutineBody(ts *TS) *Node {
	result := createNodeFromString("subroutineBody")
	result.addChild(match("{", ts))
	for ts.curTok().Contents == "var" {
		result.addChild(varDec(ts))
	}
	result.addChild(statements(ts))
	result.addChild(match("}", ts))
	return result
}

func varDec(ts *TS) *Node {
	result := createNodeFromString("varDec")
	result.addChild(match("var", ts))
	result.addChild(typeName(ts))
	result.addChild(match(IDENT, ts))
	for ts.curTok().Contents == "," {
		result.addChild(match(",", ts))
		result.addChild(match(IDENT, ts))
	}
	result.addChild(match(";", ts))
	return result
}

func statements(ts *TS) *Node {
	result := createNodeFromString("statements")
	for cur := ts.curTok(); util.Contains(statementKeywords, cur.Contents); cur = ts.curTok() {
		switch cur.Contents {
		case "let":
			result.addChild(letStatement(ts))
		case "if":
			result.addChild(ifStatement(ts))
		case "while":
			result.addChild(whileStatement(ts))
		case "do":
			result.addChild(doStatement(ts))
		case "return":
			result.addChild(returnStatement(ts))
		}
	}
	return result
}

func letStatement(ts *TS) *Node {
	result := createNodeFromString("letStatement")
	result.addChild(match("let", ts))
	result.addChild(match(IDENT, ts))
	if ts.curTok().Contents == "[" {
		result.addChild(match("[", ts))
		result.addChild(expression(ts))
		result.addChild(match("]", ts))
	}
	result.addChild(match("=", ts))
	result.addChild(expression(ts))
	result.addChild(match(";", ts))
	return result
}

func whileStatement(ts *TS) *Node {
	result := createNodeFromString("whileStatement")
	result.addChild(match("while", ts))
	result.addChild(match("(", ts))
	result.addChild(expression(ts))
	result.addChild(match(")", ts))
	result.addChild(match("{", ts))
	result.addChild(statements(ts))
	result.addChild(match("}", ts))
	return result
}

func ifStatement(ts *TS) *Node {
	result := createNodeFromString("ifStatement")
	result.addChild(match("if", ts))
	result.addChild(match("(", ts))
	result.addChild(expression(ts))
	result.addChild(match(")", ts))
	result.addChild(match("{", ts))
	result.addChild(statements(ts))
	result.addChild(match("}", ts))
	if ts.curTok().Contents == "else" {
		result.addChild(match("else", ts))
		result.addChild(match("{", ts))
		result.addChild(statements(ts))
		result.addChild(match("}", ts))
	}
	return result
}

func doStatement(ts *TS) *Node {
	result := createNodeFromString("doStatement")
	result.addChild(match("do", ts))
	result.addChild(match(IDENT, ts))
	if ts.curTok().Contents == "." {
		result.addChild(match(".", ts))
		result.addChild(match(IDENT, ts))
	}
	result.addChild(match("(", ts))
	result.addChild(expressionList(ts))
	result.addChild(match(")", ts))
	result.addChild(match(";", ts))
	return result
}

func _subroutineCallHelper(result *Node, ts *TS) *Node {
	result.addChild(match(IDENT, ts))
	if ts.curTok().Contents == "." {
		result.addChild(match(".", ts))
		result.addChild(match(IDENT, ts))
	}
	result.addChild(match("(", ts))
	result.addChild(expressionList(ts))
	result.addChild(match(")", ts))
	return result
}

func subroutineCall(ts *TS) *Node {
	result := createNodeFromString("subroutineCall")
	result = _subroutineCallHelper(result, ts)
	return result
}

func expressionList(ts *TS) *Node {
	result := createNodeFromString("expressionList")
	if ts.curTok().Contents == ")" {
		return result
	}
	result.addChild(expression(ts))
	for ts.curTok().Contents == "," {
		result.addChild(match(",", ts))
		result.addChild(expression(ts))
	}
	return result
}

func returnStatement(ts *TS) *Node {
	result := createNodeFromString("returnStatement")
	result.addChild(match("return", ts))
	if ts.curTok().Contents != ";" {
		result.addChild(expression(ts))
	}
	result.addChild(match(";", ts))
	return result
}

func expression(ts *TS) *Node {
	result := createNodeFromString("expression")
	result.addChild(term(ts))
	// will continue checking the next op if it is an operator
	for curr := ts.curTok(); util.Contains(binaryOperators, curr.Contents); curr = ts.curTok() {
		result.addChild(match(curr.Contents, ts))
		result.addChild(term(ts))
	}
	return result
}

func term(ts *TS) *Node {
	result := createNodeFromString("term")
	curr := ts.curTok()

	switch {
	case util.Contains(unaryOperators, curr.Contents):
		result.addChild(match(curr.Contents, ts))
		result.addChild(term(ts))
	case util.Contains(keywordConst, curr.Contents):
		result.addChild(match(curr.Contents, ts))
	case curr.Contents == "(":
		result.addChild(match("(", ts))
		result.addChild(expression(ts))
		result.addChild(match(")", ts))
	case curr.Kind == INT:
		result.addChild(match(INT, ts))
	case curr.Kind == STRING:
		result.addChild(match(STRING, ts))
	case curr.Kind == IDENT:
		next, err := ts.peekNextToken()
		if err != nil {
			panic(err)
		}
		if next.Contents == "[" {
			result.addChild(match(IDENT, ts))
			result.addChild(match("[", ts))
			result.addChild(expression(ts))
			result.addChild(match("]", ts))
		} else if next.Contents == "(" || next.Contents == "." {
			result = _subroutineCallHelper(result, ts)
		} else {
			result.addChild(match(IDENT, ts))
		}
    default:
        result.addChild(createNodeFromString("ERROR"))
	}

	return result
}

func op(ts *TS) *Node {
	return match(binaryOperators, ts)
}

func unaryOp(ts *TS) *Node {
	return match(unaryOperators, ts)
}

func keywordConstant(ts *TS) *Node {
	return match(keywordConst, ts)
}
