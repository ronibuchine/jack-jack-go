package tokenizer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TokenType int64

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENT
	INT
	STRING
)

const (
	NOT_A_KEYWORD = iota
)

// that can't be const for... reasons
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
	return b == ' ' || b == '\n' || b == '\t'
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


func tokenize(file string) {

	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create(strings.TrimSuffix(file, ".jack") + "T.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	reader := bufio.NewReader(input)

	var (
		tokenContents string
		tokenType     string
	)

	inComment := false

    output.WriteString("<tokens>\n")
	for {
		cur, err := reader.ReadByte()
		if err != nil {
			break
		}

		if isWhitespace(cur) {
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
        tokenType = "UNKOWN"

		// strings
		if cur == '"' {
			str, err := reader.ReadSlice('"')
			if err != nil {
				log.Fatal("Unclosed string")
			}
			tokenContents = strings.TrimSuffix(string(str), "\"")
			tokenType = "stringConstant"
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
			tokenContents = integer
			tokenType = "integerConstant"
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

			tokenContents = word
			if isKeyword(word) {
				tokenType = "keyword"
			} else {
				tokenType = "identifier"
			}
		}

		// symbols
		if validSymbol := getSymbol(cur); validSymbol != "" {
			tokenContents = validSymbol
			tokenType = "symbol"
		}

		// write token to xml
		tokenXml := fmt.Sprint("\t<" + tokenType + "> " + tokenContents + " </" + tokenType + ">\n")
		output.WriteString(tokenXml)
	}
    output.WriteString("</tokens>\n")
}

