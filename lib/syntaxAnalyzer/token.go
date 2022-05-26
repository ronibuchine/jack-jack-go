package syntaxAnalyzer

import (
	"encoding/xml"
	"errors"
)

const (
	KEYWORD = "keyword"
	SYMBOL  = "symbol"
	IDENT   = "identifier"
	INT     = "integerConstant"
	STRING  = "stringConstant"
	UNKNOWN = "UNKNOWN"
)

type Token struct {
	Kind       string
	Contents   string
	LineNumber int
}

// represents a token stream
type TS struct {
	Tokens  []Token
	counter int
	File    string
}

// this struct is just a hack for xml writing
type TokensXML struct {
	tokens []Token
}

func (ts TS) curTok() Token {
	return ts.Tokens[ts.counter]
}

// peek helper for LL(1) lookahead
func (ts TS) peekNextToken() (Token, error) {
	if ts.counter+1 >= len(ts.Tokens) {
		return Token{}, errors.New("Cannot lookahead passed the end of token stream")
	}
	return ts.Tokens[ts.counter+1], nil
}

func (t TokensXML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(xml.StartElement{Name: xml.Name{Space: "", Local: "tokens"}})
	for _, token := range t.tokens {
		e.Encode(token)
	}
	return e.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: "tokens"}})
}

func (t Token) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(xml.StartElement{Name: xml.Name{Space: "", Local: t.Kind}})
	e.EncodeToken(xml.CharData([]byte(" " + t.Contents + " ")))
	return e.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: t.Kind}})
}
