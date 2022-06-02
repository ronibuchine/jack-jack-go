package vmwriter

import (
	"bufio"
	fe "jack-jack-go/lib/syntaxAnalyzer"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	jack, _ := os.Open("redeclare.jack")
	reader := bufio.NewReader(jack)
	tokens := fe.Tokenize((reader))
	n := fe.Parse(&fe.TS{Tokens: tokens, File: "redeclare.jack"})
	j := NewJackCompiler(n, "redeclare", nil)
	j.Compile()
}
