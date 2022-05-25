package syntaxAnalyzer

import (
    "encoding/xml"
)

const (
	KEYWORD = "keyword"
	SYMBOL = "symbol"
	IDENT = "identifier"
    INT = "integerConstant"
	STRING = "stringConstant"
    UNKNOWN = "UNKNOWN"
)

// this struct is just a hack for xml writing
type TokensXML struct {
    tokens []Token
}

type Token struct {
	Kind       string
	Contents   string
	LineNumber int
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
