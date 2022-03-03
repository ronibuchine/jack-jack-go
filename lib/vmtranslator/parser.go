package vmtranslator

import (
	"bufio"
	"errors"
	"os"
	"regexp"
)

var (
	RE_ARITHMETIC    *regexp.Regexp
	RE_PUSH_POP      *regexp.Regexp
	RE_IF_LABEL_GOTO *regexp.Regexp
	RE_FUNCTION_CALL *regexp.Regexp
	RE_RETURN        *regexp.Regexp
	RE_CALL          *regexp.Regexp
)

func CompileAllRegex() {
	RE_ARITHMETIC = regexp.MustCompile(`(?m)^\s*(add|sub|eq|gt|lt|and|or|not)\s*$`)
	RE_PUSH_POP = regexp.MustCompile(`(?m)^\s*(push|pop)\s+(local|argument|this|that|constant|static|pointer|temp)\s+(\d+)\s*$`)
	RE_IF_LABEL_GOTO = regexp.MustCompile(`(?m)^\s*(if|label|goto)\s+([A-Za-z_][A-Za-z0-9_]*)\s*$`)
	RE_FUNCTION_CALL = regexp.MustCompile(`(?m)^\s*(function|call)\s+([A-Za-z_][A-Za-z0-9_]*)\s+(\d+)\s*$`)
	RE_RETURN = regexp.MustCompile(`(?m)^\s*return\s*`)
}

type C_TYPE int

const (
	C_ARITHMETIC C_TYPE = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)


type Command struct {
	cmdType C_TYPE
	arg1    string
	arg2    string
}


func parseCommand(s string) (*Command, error) {
	if cmd := RE_ARITHMETIC.FindStringSubmatch(s); cmd != nil {
		return &Command{cmdType: C_ARITHMETIC, arg1: cmd[1], arg2: ""}, nil
	}
	if cmd := RE_PUSH_POP.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "push":
			return &Command{cmdType: C_PUSH, arg1: cmd[2], arg2: cmd[3]}, nil
		case "pop":
			return &Command{cmdType: C_POP, arg1: cmd[2], arg2: cmd[3]}, nil
		}
	}
	if cmd := RE_IF_LABEL_GOTO.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "if-goto":
			return &Command{cmdType: C_IF, arg1: cmd[2], arg2: ""}, nil
		case "goto":
			return &Command{cmdType: C_GOTO, arg1: cmd[2], arg2: ""}, nil
		case "label":
			return &Command{cmdType: C_LABEL, arg1: cmd[2], arg2: ""}, nil
		}
	}
	if cmd := RE_FUNCTION_CALL.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "function":
			return &Command{cmdType: C_FUNCTION, arg1: cmd[2], arg2: cmd[3]}, nil
		case "pop":
			return &Command{cmdType: C_CALL, arg1: cmd[2], arg2: cmd[3]}, nil
		}
	}
	if cmd := RE_RETURN.FindStringSubmatch(s); cmd != nil {
		return &Command{cmdType: C_RETURN, arg1: "", arg2: ""}, nil
	}
	return nil, errors.New("command not recognized")
}

func ParseFile(file *os.File) (commands []*Command) {
    CompileAllRegex()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		command, err := parseCommand(line)
		if err != nil {
            // more error stuff??
			continue
		}
		commands = append(commands, command)
	}
	return
}
