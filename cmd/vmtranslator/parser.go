package vmtranslator

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
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
	RE_ARITHMETIC = regexp.MustCompile("(?m)^\\s*(add|sub|eq|gt|lt|and|or|not)\\s*$")
	RE_PUSH_POP = regexp.MustCompile("(?m)^\\s*(push|pop)\\s+(local|argument|this|that|constant|static|pointer|temp)\\s+(\\d+)\\s*$")
	RE_IF_LABEL_GOTO = regexp.MustCompile("(?m)^\\s*(if|label|goto)\\s+([A-Za-z_][A-Za-z0-9_]*)\\s*$")
	RE_FUNCTION_CALL = regexp.MustCompile("(?m)^\\s*(function|call)\\s+([A-Za-z_][A-Za-z0-9_]*)\\s+(\\d+)\\s*$")
	RE_RETURN = regexp.MustCompile("(?m)^\\s*return\\s*")
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

type Command interface {
	cmdType() C_TYPE
	arg1() (string, error)
	arg2() (int, error)
}

type Command_t struct {
	_cmdType C_TYPE
	_arg1    string
	_arg2    int
}

func (c Command_t) cmdType() C_TYPE {
	return c._cmdType
}

func (c Command_t) arg1() (string, error) {
	if c._cmdType == C_RETURN {
		return "", errors.New("return command does not have arg 1")
	} else {
		return c._arg1, nil
	}
}

func (c Command_t) arg2() (int, error) {
	if c._cmdType == C_PUSH || c._cmdType == C_POP || c._cmdType == C_FUNCTION || c._cmdType == C_CALL {
		return c._arg2, nil
	} else {
		return -1, errors.New("The current command does not have arg 2")
	}
}

func parseCommand(s string) (Command, error) {
	if cmd := RE_ARITHMETIC.FindStringSubmatch(s); cmd != nil {
		return Command_t{_cmdType: C_ARITHMETIC, _arg1: cmd[1], _arg2: -1}, nil
	}
	if cmd := RE_PUSH_POP.FindStringSubmatch(s); cmd != nil {
		arg2, err := strconv.Atoi(cmd[3])
		if err != nil {
			// error??
		}
		switch cmd[1] {
		case "push":
			return Command_t{_cmdType: C_PUSH, _arg1: cmd[2], _arg2: arg2}, nil
		case "pop":
			return Command_t{_cmdType: C_POP, _arg1: cmd[2], _arg2: arg2}, nil
		}
	}
	if cmd := RE_IF_LABEL_GOTO.FindStringSubmatch(s); cmd != nil {
		switch cmd[1] {
		case "if-goto":
			return Command_t{_cmdType: C_IF, _arg1: cmd[2], _arg2: -1}, nil
		case "goto":
			return Command_t{_cmdType: C_GOTO, _arg1: cmd[2], _arg2: -1}, nil
		case "label":
			return Command_t{_cmdType: C_LABEL, _arg1: cmd[2], _arg2: -1}, nil
		}
	}
	if cmd := RE_FUNCTION_CALL.FindStringSubmatch(s); cmd != nil {
		arg2, err := strconv.Atoi(cmd[3])
		if err != nil {
			// error stuff ??
		}
		switch cmd[1] {
		case "function":
			return Command_t{_cmdType: C_FUNCTION, _arg1: cmd[2], _arg2: arg2}, nil
		case "pop":
			return Command_t{_cmdType: C_CALL, _arg1: cmd[2], _arg2: arg2}, nil
		}
	}
	if cmd := RE_RETURN.FindStringSubmatch(s); cmd != nil {
		return Command_t{_cmdType: C_RETURN, _arg1: "", _arg2: -1}, nil
	}
	return nil, errors.New("Command not recognized")
}


func ParseFile(file *os.File) (commands []Command) {
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
