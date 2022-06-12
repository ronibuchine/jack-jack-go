package vmwriter

import (
	"errors"
	"fmt"
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
			"static":   0,
			"this":     0,
			"argument": 0,
			"local":    0,
		}}
}

func formatError(ln int, err error) error {
    return fmt.Errorf("On line: %d: %v", ln, err)
}

// returns error if cannot add variable to symbol table
func (st *SymbolTable) Add(kind string, vType string, name string) error {
	id := st.counts[kind]
	st.counts[kind]++
	if _, exists := st.entries[name]; exists {
		return errors.New("Variable " + name + " redeclared here")
	}
	st.entries[name] = TableEntry{kind, vType, id}
	return nil
}

func (st *SymbolTable) Clear() {
	*st = *newSymbolTable()
}

// should only be passed static, field, arg, or local
func (st *SymbolTable) Count(kind string) int {
	return st.counts[kind]
}

func (st *SymbolTable) Find(name string) (TableEntry, error) {
	if entry, ok := st.entries[name]; ok {
		return entry, nil
	} else {
		return TableEntry{}, errors.New("Could not locate " + name + " in symbol table")
	}
}
