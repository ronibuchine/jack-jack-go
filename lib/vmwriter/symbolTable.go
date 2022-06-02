package vmwriter

import (
	"errors"
	fe "jack-jack-go/lib/syntaxAnalyzer"
)

type TableEntry struct {
	kind  string
	vType string
	id    int
}

type SymbolTable struct {
	entries map[string]TableEntry
	counts  map[string]int
}

func newSymbolTable() *SymbolTable {
	return &SymbolTable{
		entries: make(map[string]TableEntry),
		counts: map[string]int{
			"static": 0,
			"field":  0,
			"arg":    0,
			"var":    0,
		}}
}

// create a symbol table from all the class variable declarations
func ClassTable(varDecs []*fe.Node) (*SymbolTable, error) {
	st := newSymbolTable()
	var kind, vType, name string
	for _, varDec := range varDecs {
		kind = varDec.Children[0].Token.Contents
		vType = varDec.Children[1].Token.Contents
		for i := 2; i < len(varDec.Children); i += 2 {
			name = varDec.Children[i].Token.Contents
			err := st.Add(name, kind, vType)
			if err != nil {
				return nil, err
			}
		}
	}
	return st, nil
}

// expects a node of kind subroutineDec
func LocalTable(subroutine *fe.Node) (*SymbolTable, error) {
	lst := newSymbolTable()
	var name, vType string
	if functionKind := subroutine.Children[0].Token.Kind; functionKind == "constructor" || functionKind == "method" {
		vType = subroutine.Children[1].Token.Contents
		if vType == "void" && functionKind == "constructor" {
			return nil, errors.New("void constructor makes no sense")
		}
		lst.Add("this", "arg", vType)
	}
	var params *fe.Node
	for _, child := range subroutine.Children {
		if child.Token.Kind == "parameterList" {
			params = child
			break
		}
	}
	for i := 0; i < len(params.Children); i += 3 {
		vType = params.Children[i].Token.Contents
		name = params.Children[i+1].Token.Contents
		err := lst.Add(name, "arg", vType)
		if err != nil {
			return nil, err
		}
	}
	return lst, nil
}

// returns error if cannot add variable to symbol table
func (st *SymbolTable) Add(name string, kind string, vType string) error {
	id := st.counts[kind]
	st.counts[kind]++
	if _, exists := st.entries[name]; exists {
		return errors.New("Variable " + name + " redeclared here")
	}
	st.entries[name] = TableEntry{kind, vType, id}
	return nil
}

// should only be passed static, field, arg, or local
func (st *SymbolTable) Count(kind string) int {
	return st.counts[kind]
}

func (st *SymbolTable) KindOf(name string) (string, error) {
	if entry, ok := st.entries[name]; ok {
		return entry.kind, nil
	} else {
		return "", errors.New("Could not locate " + name + " in symbol table")
	}
}

func (st *SymbolTable) TypeOf(name string) (string, error) {
	if entry, ok := st.entries[name]; ok {
		return entry.vType, nil
	} else {
		return "", errors.New("Could not locate " + name + " in symbol table")
	}
}

func (st *SymbolTable) IndexOf(name string) (int, error) {
	if entry, ok := st.entries[name]; ok {
		return entry.id, nil
	} else {
		return 0, errors.New("Could not locate " + name + " in symbol table")
	}
}
