package vmwriter

import (
    fe "jack-jack-go/lib/syntaxAnalyzer"
)

type SymbolTable struct {
}

func ClassTable(vars []*fe.Node) *SymbolTable {
    return nil
}

func LocalTable(vars []*fe.Node) *SymbolTable {
    return nil
}

func (s SymbolTable) Add(name string, varType string, kind string) error {
    return nil
}

func (s SymbolTable) Count(kind string) int {
    return 0
}

func (s SymbolTable) KindOf(name string) (string, error) {
    return "",nil
}

func (s SymbolTable) TypeOf(name string) (string, error) {
    return "",nil
}

// running index should be kept seperate for each kind of variable
func (s SymbolTable) IndexOf(name string) (string, error) {
    return "",nil
}
