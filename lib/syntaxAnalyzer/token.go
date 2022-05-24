package syntaxAnalyzer

const (
	KEYWORD = "keyword"
	SYMBOL = "symbol"
	IDENT = "identifier"
    INT = "integerConstant"
	STRING = "stringConstant"
    UNKNOWN = "UNKNOWN"
)

type Token struct {
	Kind       string
	Contents   string
	LineNumber int
}
