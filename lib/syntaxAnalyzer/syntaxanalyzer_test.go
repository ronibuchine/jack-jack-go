package syntaxAnalyzer

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	var b bytes.Buffer
	jack, _ := os.Open("../../test.jack")
	reader := bufio.NewReader(jack)
	tokens := Tokenize((reader))
	NodeToXML(Parse(&TS{Tokens: tokens, File: "test.jack"}), &b)
	fmt.Print(b.String())
}

func TestTokenize(t *testing.T) {
	var b bytes.Buffer
	jack, _ := os.Open("test.jack")
	reader := bufio.NewReader(jack)
	TokenToXML(Tokenize(reader), &b)
	fmt.Print(b.String())
}
