package compiler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// aliasing structs for xml objects
type TokensXml struct {
	XMLName xml.Name   `xml:"tokens"`
	Tokens  []TokenXml `xml:"token"`
}

type TokenXml struct {
	XMLName xml.Name `xml:"token"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}

// Node struct for parsing and recursive descent
type Node struct {
	/* name     string
	contents string // will be the same as name for non-terminals */
	token    Token
	children []*Node
}

// cant be const for some reason
var operators []string = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func createNodeFromToken(token Token) *Node {
	return &Node{token, []*Node{}}
}

func createNodeFromString(name string) *Node {
	return &Node{Token{Kind: name}, []*Node{}}
}

func (parent *Node) addChild(child *Node) {
	parent.children = append(parent.children, child)
}

func (n *Node) name() string {
	return n.token.Kind
}

// globals for matching
var (
	TokenStream  []Token
	tokenCounter int
)

func curTok() Token {
	return TokenStream[tokenCounter]
}

// peek helper for LL(1) lookahead
func peekNextToken() Token {
	return TokenStream[tokenCounter+1]
}

// Takes an ordered token stream which is a map[string]string and parses it and returns the XML tree
func BuildXML(tokenStream map[string]string) *os.File {
	output, err := os.Create("output.xml")
	if err != nil {
		log.Fatal("Failed to create an ouput file.")
	}
	defer output.Close()
	// TODO: implement

	return output
}

func GetTokens(jackFile string) {
	TokenStream = Tokenize(jackFile)
	tokenCounter = 0
}

func _matchSingle(token string) (*Node, error) {
	// if we match a ident, int or string, we DONT care about the contents
	// if we match a symbol or keyword, we DO care about contents
	curr := curTok()
	if ((token == IDENT || token == INT || token == STRING) && (token == curr.Kind)) ||
		((token == curr.Contents) && (curr.Kind == KEYWORD || curr.Kind == SYMBOL)) {
		res := createNodeFromToken(curr)
		tokenCounter++
		return res, nil
	}
	// TODO: add error handling. We might just panic and die here
	// error(fmt.Sprint("Expected token %s before ", token))
	return createNodeFromString("ERROR"), errors.New(fmt.Sprint("Failed to match %s", token))
}

// global function used for token parsing and matching
// can either pass in a string or []string. If no matches, then an error will
// occur, if at least one matches then the first match will be returned
func match(token interface{}) (result *Node) {
	if tokenCounter >= len(TokenStream) {
		log.Fatal("end of token stream")
	}

	if t, ok := token.(string); ok {
		if res, err := _matchSingle(t); err == nil {
			result = res
		} else {
			parseError(t)
			result = createNodeFromString("ERROR")
		}
	} else if tokens, ok := token.([]string); ok {
		for _, t := range tokens {
			if res, err := _matchSingle(t); err == nil {
				result = res
			}
		}
		parseError(strings.Join(tokens, ", "))
		result = createNodeFromString("ERROR")
	}
	return result
}

func parseError(expected string) {
	curr := curTok()
	fmt.Sprint(fmt.Sprint("ERROR line %d: Expected token(s) `%s` before %s %s\n", curr.LineNumber, expected, curr.Kind, curr.Contents))
}

// functions for grammar
func class() *Node {
	result := createNodeFromString("class")
	result.addChild(match("class"))
	result.addChild(match(IDENT))
	result.addChild(match("{"))
	curr := curTok()
	for curr.Kind == KEYWORD &&
		(curr.Contents == "static" || curr.Contents == "field") {
		result.addChild(classVarDec())
		curr = curTok()
	}
	for curr.Kind == KEYWORD &&
		(curr.Contents == "constructor" || curr.Contents == "function" || curr.Contents == "method") {
		result.addChild(subroutineDec())
		curr = curTok()
	}
	result.addChild(match("}"))
	return result
}

func classVarDec() *Node {
	result := createNodeFromString("classVarDec")
	result.addChild(match([]string{"static", "field"}))
	result.addChild(typeName())
	result.addChild(match(IDENT))
	for curTok().Contents == "," {
		result.addChild(match(","))
		result.addChild(match(IDENT))
	}
	result.addChild(match(";"))
	return result
}

func typeName() *Node {
	result := match([]string{"int", "char", "boolean", IDENT})
	return result
}

func subroutineDec() *Node {
	result := createNodeFromString("subroutineDec")
	result.addChild(match([]string{"constructor", "function", "method"}))
	if curTok().Contents == "void" {
		result.addChild(match("void"))
	} else {
		result.addChild(typeName())
	}
	result.addChild(match(IDENT))
	result.addChild(match("("))
	result.addChild(parameterList())
	result.addChild(match(")"))
	result.addChild(subroutineBody())
	return result
}

func parameterList() *Node {
	result := createNodeFromString("parameterList")
	if curTok().Contents == ")" {
		return result
	}
	result.addChild(typeName())
	result.addChild(match(IDENT))
	for curTok().Contents == "," {
		result.addChild(match(","))
		result.addChild(typeName())
		result.addChild(match(IDENT))
	}
	return result
}

func subroutineBody() *Node {
	result := createNodeFromString("subroutineBody")
	result.addChild(match("{"))
	for curTok().Contents == "var" {
		result.addChild(varDec())
	}
	result.addChild(statements())
	result.addChild(match("}"))
	return result
}

func varDec() *Node {
	result := createNodeFromString("varDec")
	result.addChild(match("var"))
	result.addChild(typeName())
	result.addChild(match(IDENT))
	for curTok().Contents == "," {
		result.addChild(match(","))
		result.addChild(match(IDENT))
	}
	result.addChild(match(";"))
	return result
}

func statements() *Node {
	panic("unimplemented")
}

func statement() *Node {
	result := createNodeFromString("statement")
	switch curTok().Contents {
	case "let":
		result.addChild(letStatement())
	case "if":
		result.addChild(ifStatement())
	case "while":
		result.addChild(whileStatement())
	case "do":
		result.addChild(doStatement())
	case "return":
		result.addChild(returnStatement())
	}
	return result
}

func letStatement() *Node {
	result := createNodeFromString("letStatement")
	result.addChild(match("let"))
	result.addChild(match(IDENT))
	if curTok().Contents == "[" {
		result.addChild(match("["))
		result.addChild(match(expression()))
		result.addChild(match("]"))
	}
	result.addChild(match("="))
	result.addChild(match(expression()))
	result.addChild(match(";"))
	return result
}

func whileStatement() *Node {
	result := createNodeFromString()
	result.addChild(match("whileStatement"))
	result.addChild(match("("))
	result.addChild(match(expression()))
	result.addChild(match(")"))
	result.addChild(match("{"))
	result.addChild(match(statements()))
	result.addChild(match("}"))
}

func ifStatement() *Node {
	result := createNodeFromString("ifStatement")
	result.addChild(match("if"))
	result.addChild(match("("))
	result.addChild(match(expression()))
	result.addChild(match(")"))
	result.addChild(match("{"))
	result.addChild(match(statements()))
	result.addChild(match("}"))
	if curTok().Contents == "else" {
		result.addChild(match("else"))
		result.addChild(match("{"))
		result.addChild(match(statements()))
		result.addChild(match("}"))
	}
	return result
}

func doStatement() *Node {
	result := createNodeFromString("doStatement")
	result.addChild(match("do"))
	result.addChild(match(subroutineCall()))
	result.addChild(match(";"))
	return result
}

func subroutineCall() *Node {
	result := createNodeFromString("subroutineCall")
	result.addChild(match(IDENT))
	if curTok().Contents == "." {
		result.addChild(match("."))
	}
	result.addChild(match("("))
	result.addChild(match(expressionList()))
	result.addChild(match(")"))
	return result
}

func expressionList() *Node {
	panic("unimplemented")
}

func returnStatement() *Node {
	result := createNodeFromString("returnStatement")
	result.addChild(match("return"))
	if curTok().Contents != "}" {
		result.addChild(match(expression()))
	}
	result.addChild(match(";"))
	return result
}

// helper function to check existence in a collection, for some reason this doesnt exist in the go stdlib...
// if you want to use for other types just add to the generic parameter list
// doesn't return a bool, returns the item itself otherwise returns nil
func _contains[T string | int](collection []T, item T) (T, error) {
	for _, value := range collection {
		if item == value {
			return value, nil
		}
	}
	return item, errors.New("The collection does not contain the given item")
}

func expression() *Node {
	result := createNodeFromString("expression")
	result.addChild(match(term()))
	// will continue checking the next op if it is an operator
	for op, err := _contains(operators, curTok().Contents); err != nil; {
		result.addChild(match(op))
		result.addChild(match(term()))
	}
	return result
}

func term() *Node {
	result := createNodeFromString("term")
	// finish implementation
	return result
}

func constant() *Node {
	return match(INT)
}

func op() *Node {
	return match([]string{"+", "-", "*", "/", "&", "|", "<", ">", "="})
}

func unaryOp() *Node {
	return match([]string{"~", "-"})
}

func keywordConstant() *Node {
	return match([]string{"true", "false", "null", "this"})
}
