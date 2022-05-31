package vmwriter

import (
	fe "jack-jack-go/lib/syntaxAnalyzer"
)

type SymbolTable struct {
	name        string
	kind        string
	varType     string
	argCount    int
	localsCount int
	staticCount int
	fieldCount  int
}

// create a symbol table from all the class variable declarations
func ClassTable(vars []*fe.Node) *SymbolTable {

	return nil
}

// expects a node of kind subroutineDec
func LocalTable(subroutine *fe.Node) *SymbolTable {
	if subroutine.Children[0].Token.Kind == "constructor" { // euggh
		// add this as arg 0
	}
	var params *fe.Node
	for _, child := range subroutine.Children {
		if child.Token.Kind == "parameterList" {
			params = child
		}
	}
	lst := &SymbolTable{}
	numChildren := len(params.Children)
	if numChildren == 2 {
		// add type and identifier to lst
		return lst
	} else { // children > 2
		// for each param in list, add to symbol table
		for i := 0; i < numChildren; i += 3 {
			// type is params.Children[i]
			// identifier is params.Children[i+1]
			// skip the comma
		}
	}
	return lst
}

func (s SymbolTable) Add(name string, varType string, kind string) error {
	return nil
}

func (s SymbolTable) Count(kind string) int {
	return 0
}

func (s SymbolTable) KindOf(name string) (string, error) {
	return "", nil
}

func (s SymbolTable) TypeOf(name string) (string, error) {
	return "", nil
}

// running index should be kept seperate for each kind of variable
func (s SymbolTable) IndexOf(name string) (string, error) {
	return "", nil
}
