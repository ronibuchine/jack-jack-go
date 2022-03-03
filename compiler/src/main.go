package backend

import (
	"errors"
	"os"
)

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

type ARITHMETIC_TYPE string

const (
	ADD ARITHMETIC_TYPE = "add"
	SUB ARITHMETIC_TYPE = "sub"
	NEG ARITHMETIC_TYPE = "neg"
	EQ  ARITHMETIC_TYPE = "eq"
	GT  ARITHMETIC_TYPE = "gt"
	LT  ARITHMETIC_TYPE = "lt"
	AND ARITHMETIC_TYPE = "and"
	OR  ARITHMETIC_TYPE = "or"
	NOT ARITHMETIC_TYPE = "not"
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


}

func parser(file *os.File) ([]Command, error) {

}
