package syntaxAnalyzer

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var KEYWORDS_LIST []string = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

var SYMBOL_LIST = map[byte]string{
	'{': "{",
	'}': "}",
	'(': "(",
	')': ")",
	'[': "[",
	']': "]",
	'.': ".",
	',': ",",
	';': ";",
	'+': "+",
	'-': "-",
	'*': "*",
	'/': "/",
	'&': "&amp",
	'|': "|",
	'<': "&lt",
	'>': "&gt",
	'=': "=",
	'~': "~",
}

func peekByte(r *bufio.Reader) byte {
	b, err := r.Peek(1)
	if err != nil {
		return 0
	}
	return b[0]
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isWordStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isKeyword(word string) bool {
	for _, k := range KEYWORDS_LIST {
		if k == word {
			return true
		}
	}
	return false
}

func getSymbol(b byte) string {
	if symbol, ok := SYMBOL_LIST[b]; ok {
		return symbol
	}
	return ""
}

func writeXMLHeader(output *os.File) error {
	if _, err := output.WriteString(`<?xml version="1.0" encoding="UTF-8" ?>` + "\n"); err != nil {
		return err
	}
	return nil
}

func TokenToXML(tokens []Token, w io.Writer) error {
	bytes, err := xml.MarshalIndent(TokensXML{tokens}, "", "    ")
	bytes = []byte(xml.Header + string(bytes))
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

func Tokenize(file string) []Token {

	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create(strings.TrimSuffix(file, ".jack") + "_tokenized.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	reader := bufio.NewReader(input)

	var curToken Token

	inComment := false
	if err := writeXMLHeader(output); err != nil {
		log.Fatal("Failed to write the XML header to the output file.\n")
	}
	output.WriteString("<tokens>\n")

	var tokenStream []Token

	lineNumber := 1
	for {
		cur, err := reader.ReadByte()
		if err != nil {
			break
		}

		if isWhitespace(cur) {
			continue
		}

		if cur == '\n' || cur == '\r' {
			lineNumber += 1
			continue
		}

		next := peekByte(reader)

		// end of comment block
		if cur == '*' && next == '/' {
			if inComment {
				reader.Discard(1)
				inComment = false
			} else {
				log.Fatal("End of comment encountered outside of comment block")
			}
			continue
		}

		// skip forward until end of comment block
		if inComment {
			continue
		}

		if cur == '/' && next == '/' {
			_, err := reader.ReadSlice('\n')
			lineNumber += 1
			if err != nil {
				log.Fatal("Unclosed line comment")
			}
			continue
		}

		// start of comment block
		if cur == '/' && next == '*' {
			inComment = true
			reader.Discard(1)
			continue
		}

		// now figure out what the next token is
		curToken.Kind = UNKNOWN
		curToken.Contents = string(cur)

		// strings
		if cur == '"' {
			str, err := reader.ReadSlice('"')
			newLineCount := strings.Count(string(str), "\n")
			lineNumber += newLineCount
			if err != nil {
				log.Fatal("Unclosed string")
			}
			curToken.Contents = strings.TrimSuffix(string(str), "\"")
			curToken.Kind = STRING
		}

		// digits
		if isDigit(cur) {
			var integer string
			for c := cur; isDigit(c); c, err = reader.ReadByte() {
				if err != nil {
					log.Fatal("Ended file in middle of number")
				}
				integer += string(c)
			}
			reader.UnreadByte()

			value, err := strconv.Atoi(integer)
			if err != nil || value < 0 || value > 32767 {
				log.Fatal("Number is invalid")
			}
			curToken.Contents = integer
			curToken.Kind = INT
		}

		// idents and keywords
		if isWordStart(cur) {
			var word string
			for c := cur; isWordStart(c) || isDigit(c); c, err = reader.ReadByte() {
				if err != nil {
					log.Fatal("Ended file in middle of word") // wut is this error message
				}
				word += string(c)
			}
			reader.UnreadByte()

			curToken.Contents = word
			if isKeyword(word) {
				curToken.Kind = KEYWORD
			} else {
				curToken.Kind = IDENT
			}
		}

		// symbols
		if validSymbol := getSymbol(cur); validSymbol != "" {
			curToken.Contents = validSymbol
			curToken.Kind = SYMBOL
		}

		// write token to xml
		tokenXml := fmt.Sprintf("<%s> %s </%s>\n", curToken.Kind, curToken.Contents, curToken.Kind)
		output.WriteString(tokenXml)

		curToken.LineNumber = lineNumber
		tokenStream = append(tokenStream, curToken)
	}

	output.WriteString("</tokens>\n")
	return tokenStream
}
