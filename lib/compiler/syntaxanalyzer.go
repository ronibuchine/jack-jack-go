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
type TokensXML struct {
	XMLName xml.Name   `xml:"tokens"`
	Tokens  []TokenXML `xml:"token"`
}

type TokenXML struct {
	XMLName xml.Name `xml:"token"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}

// Node struct for parsing and recursive descent
type Node struct {
	token    Token
	children []*Node
}

func createNodeFromToken(token Token) *Node {
	return &Node{token, []*Node{}}
}

func createNodeFromString(name string) *Node {
	return &Node{Token{Kind: name}, []*Node{}}
}

func (parent *Node) addChild(child *Node) {
	parent.children = append(parent.children, child)
}

// globals for matching
var (
	TokenStream  []Token
	tokenCounter int
)

func curTok() Token {
	return TokenStream[tokenCounter]
}

// cant be const for some reason
var binaryOperators []string = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
var unaryOperators []string = []string{"~", "-"}
var keywordConst []string = []string{"true", "false", "null", "this"}
var functionDecs []string = []string{"function", "constructor", "method"}
var classVars []string = []string{"static", "field"}

// peek helper for LL(1) lookahead
func peekNextToken() (Token, error) {
	if tokenCounter+1 >= len(TokenStream) {
		return Token{}, errors.New("Cannot lookahead passed the end of token stream")
	}
	return TokenStream[tokenCounter+1], nil
}

// Takes an ordered token stream which is a map[string]string and parses it and returns the XML tree
func BuildXML() {
	output, err := os.Create("output.xml")
	if err != nil {
		log.Fatal("Failed to create an ouput file.")
	}
	defer output.Close()
	if TokenStream[0].Contents != "class" {
		log.Fatal("Jack file must be contained in a class object")
	}
	rootNode := class()
	bytes, err := xml.MarshalIndent(rootNode, "", "    ")
	if err != nil {
		log.Fatal("Failed to build the XML fromt the root class Node")
	}
	output.Write(bytes)
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
	for _contains(classVars, curr.Contents) {
		result.addChild(classVarDec())
		curr = curTok()
	}
	for _contains(functionDecs, curr.Contents) {
		result.addChild(subroutineDec())
		curr = curTok()
	}
	result.addChild(match("}"))
	return result
}

func classVarDec() *Node {
	result := createNodeFromString("classVarDec")
	result.addChild(match(classVars))
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
	result.addChild(match(functionDecs))
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
	result := createNodeFromString("statements")
	for _contains([]string{"let", "do", "if", "while", "return"}, curTok().Contents) {
		result.addChild(statement())
	}
	return result
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
	result := createNodeFromString("whileStatement")
	result.addChild(match("("))
	result.addChild(match(expression()))
	result.addChild(match(")"))
	result.addChild(match("{"))
	result.addChild(match(statements()))
	result.addChild(match("}"))
	return result
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
        result.addChild(match(IDENT))
	}
	result.addChild(match("("))
	result.addChild(match(expressionList()))
	result.addChild(match(")"))
	return result
}

func expressionList() *Node {
	result := createNodeFromString("expressionList")
	if curTok().Contents == ")" {
		return result
	}
	result.addChild(match(expression()))
	for curTok().Contents == "," {
		result.addChild(match(","))
		result.addChild(match(expression()))
	}
	return result
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
func _contains[T string | int](collection []T, item T) bool {
	for _, value := range collection {
		if item == value {
			return true
		}
	}
	return false
}

func expression() *Node {
	result := createNodeFromString("expression")
	result.addChild(match(term()))
	// will continue checking the next op if it is an operator
	curr := curTok()
	for _contains(binaryOperators, curr.Contents) {
		result.addChild(match(curr.Contents))
		result.addChild(match(term()))
	}
	return result
}

func term() *Node {
	result := createNodeFromString("term")
	curr := curTok()

	switch {
	case _contains(binaryOperators, curr.Contents):
		result.addChild(match(curr.Contents))
		result.addChild(term())
	case _contains(keywordConst, curr.Contents):
		result.addChild(match(curr.Contents))
	case curr.Contents == "(":
		result.addChild(match("("))
		result.addChild(expression())
		result.addChild(match(")"))
	case curr.Kind == INT:
		result.addChild(match(INT))
	case curr.Kind == STRING:
		result.addChild(match(STRING))
	case curr.Kind == IDENT:
		next, err := peekNextToken()
		if err != nil {
			panic(err)
		}
		if next.Contents == "[" {
			result.addChild(match(IDENT))
			result.addChild(match("["))
			result.addChild(expression())
			result.addChild(match("]"))
		} else if next.Contents == "(" || next.Contents == "." {
			result.addChild(subroutineCall())
		} else {
			result.addChild(match(IDENT))
		}
	}
	return result
}

func op() *Node {
	return match(binaryOperators)
}

func unaryOp() *Node {
	return match(unaryOperators)
}

func keywordConstant() *Node {
	return match(keywordConst)
}
