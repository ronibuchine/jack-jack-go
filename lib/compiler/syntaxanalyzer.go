package compiler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// aliasing structs for xml objects
type Tokens struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token  `xml:"tokens"`
}

type Token struct {
	XMLName xml.Name `xml:"token"`
	Type    string   `xml:"type,attr"`
	Token   string   `xml:"token"`
}

// Node struct for parsing and recursive descent
type Node struct {
	name     string
	children []Node
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

// Reads in the token stream form the XML file, returns as a list of tokens
func ReadStream(tokenStream string) [][]string {
	input, err := os.Open(tokenStream)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	byteValue, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatal("Error reading the bytes of the XML file.")
	}

	var tokens Tokens
	xml.Unmarshal(byteValue, &tokens)
	for _, token := range tokens.Tokens {
		fmt.Println("Type: " + token.Type)
		fmt.Println("XMLName: " + token.XMLName.Local)
	}
	return getTokenStrings(tokens.Tokens)
}

// returns an array of tuples, first value in the tuple is the token, second value is the token type
func getTokenStrings(tokens []Token) [][]string {
	var tokenStrings [][]string
	for i := 0; i < len(tokens); i++ {
		tokenStrings[i][0] = tokens[i].Token
		tokenStrings[i][1] = tokens[i].Type
	}
	return tokenStrings
}

// global variable used for token parsing and matching
var current int = 0

func match(tokens []Token, token string) Node {
	if current >= len(tokens) {
		log.Fatal("end of token stream")
	}
	if token == tokens[current].XMLName.Local {
		return Node{token, []Node{}}
	} else {
		return Node{"", []Node{}}
	}
}
