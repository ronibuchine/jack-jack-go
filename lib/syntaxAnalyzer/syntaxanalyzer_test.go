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
	jack, _ := os.Open("../../tests/10/Square/Square.jack")
	reader := bufio.NewReader(jack)
	NodeToXML(Parse(Tokenize(reader)), &b)
	fmt.Print(b.String())
}

func TestTokenize(t *testing.T) {
	var b bytes.Buffer
	jack, _ := os.Open("test.jack")
	reader := bufio.NewReader(jack)
	TokenToXML(Tokenize(reader), &b)
	fmt.Print(b.String())
}
