package syntaxAnalyzer

import (
	"bufio"
	"encoding/xml"
	"io"
	"jack-jack-go/lib/util"
	"log"
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

var SYMBOL_LIST = []byte{
	'{',
	'}',
	'(',
	')',
	'[',
	']',
	'.',
	',',
	';',
	'+',
	'-',
	'*',
	'/',
	'&',
	'|',
	'<',
	'>',
	'=',
	'~',
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
	return b == ' ' || b == '\t' || b == '\r'
}

func TokenToXML(tokens []Token, w io.Writer) error {
	bytes, err := xml.MarshalIndent(TokensXML{tokens}, "", "  ")
	if err != nil {
		return err
	}
	// bytes = []byte(xml.Header + string(bytes))
	w.Write(bytes)
	return nil
}

func Tokenize(reader *bufio.Reader) []Token {

	inComment := false
	var curToken Token
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

		if cur == '\n' {
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
			if util.Contains(KEYWORDS_LIST, word) {
				curToken.Kind = KEYWORD
			} else {
				curToken.Kind = IDENT
			}
		}

		// symbols
		if util.Contains(SYMBOL_LIST, cur) {
			curToken.Contents = string(cur)
			curToken.Kind = SYMBOL
		}

		curToken.LineNumber = lineNumber
		tokenStream = append(tokenStream, curToken)
	}

    return tokenStream
}
