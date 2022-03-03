package vmtranslator

import "log"

type AC_TYPE int

const (
	AC_OP AC_TYPE = iota
	AC_RELOP
	AC_VAR
	AC_JUMP
)

type AsmCommand_t struct{}

type AsmCommand interface{}

func TranslateCommand(cmd Command) (asmCommand string, err error) {
	log.Fatal("unimplemented")
	return
}
