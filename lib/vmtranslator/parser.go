package vmtranslator

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

var (
	ReArithmetic   *regexp.Regexp
	RePushPop      *regexp.Regexp
	ReIfLabelGoto  *regexp.Regexp
	ReFunctionCall *regexp.Regexp
	ReReturn       *regexp.Regexp
	ReComment      *regexp.Regexp
)

func CompileAllRegex() {
	// MustCompile takes care of lazy compilation of regex
	ReArithmetic = regexp.MustCompile(`(?m)^\s*(add|sub|neg|eq|gt|lt|and|or|not)\s*$`)
	RePushPop = regexp.MustCompile(`(?m)^\s*(push|pop)\s+(local|argument|this|that|constant|static|pointer|temp)\s+(\d+)\s*$`)
	ReIfLabelGoto = regexp.MustCompile(`(?m)^\s*(label|goto|if-goto)\s+([A-Za-z_][A-Za-z0-9_]*)\s*$`)
	ReFunctionCall = regexp.MustCompile(`(?m)^\s*(function|call)\s+([A-Za-z_][A-Za-z0-9_]*)\s+(\d+)\s*$`)
	ReReturn = regexp.MustCompile(`(?m)^\s*return\s*$`)
	ReComment = regexp.MustCompile(`(?m)(\/\/.*$)`)
}

type CType int

const (
	_ CType = iota
	CArithmetic
	CPush
	CPop
	CLabel
	CGoto
	CIf
	CFunction
	CReturn
	CCall
)

type Command struct {
	cmdType CType
	arg1    string
	arg2    string
}

func parseCommand(s string, translationUnit string) (*Command, error) {

	s = ReComment.ReplaceAllLiteralString(s, "")

	if cmd := ReArithmetic.FindStringSubmatch(s); cmd != nil {
		return &Command{cmdType: CArithmetic, arg1: cmd[1], arg2: ""}, nil
	}
	if cmd := RePushPop.FindStringSubmatch(s); cmd != nil {
		arg2 := cmd[3]
		if cmd[2] == "static" {
			arg2 = translationUnit + "." + cmd[3]
		}
		switch cmd[1] {
		case "push":
			return &Command{cmdType: CPush, arg1: cmd[2], arg2: arg2}, nil
		case "pop":
			return &Command{cmdType: CPop, arg1: cmd[2], arg2: arg2}, nil
		}
	}
	if cmd := ReIfLabelGoto.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "if-goto":
			return &Command{cmdType: CIf, arg1: cmd[2], arg2: ""}, nil
		case "goto":
			return &Command{cmdType: CGoto, arg1: cmd[2], arg2: ""}, nil
		case "label":
			return &Command{cmdType: CLabel, arg1: cmd[2], arg2: ""}, nil
		}
	}
	if cmd := ReFunctionCall.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "function":
			return &Command{cmdType: CFunction, arg1: cmd[2], arg2: cmd[3]}, nil
		case "pop":
			return &Command{cmdType: CCall, arg1: cmd[2], arg2: cmd[3]}, nil
		}
	}
	if cmd := ReReturn.FindStringSubmatch(s); cmd != nil {
		return &Command{cmdType: CReturn, arg1: "", arg2: ""}, nil
	}
	return nil, errors.New("command not recognized")
}

func ParseFile(file io.Reader, translationUnit string) (commands []*Command) {
	CompileAllRegex()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		command, err := parseCommand(line, translationUnit)
		if err != nil {
			// more error stuff??
			continue
		}
		commands = append(commands, command)
	}
	return
}

func (cmd Command) ToString() string {
	return "\nCommand Type: " + string(rune(cmd.cmdType)) + "\narg1: " + cmd.arg1 + "\narg2: " + cmd.arg2
}
