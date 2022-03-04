package vmtranslator

import (
	"strings"
)

func pushPopToHack(command Command) (hack string) {
	if segment := command.cmdType; segment == C_PUSH {

		// Place VM argument data in D register
		hack = vmArgumentAddressToAD(command) + "D=M\n"

		hack += pushFromD()

	} else if segment == C_POP {
		// Place VM argument address in M[13] (temp)
		hack += vmArgumentAddressToAD(command) + "@R13\nM=D\n"
		// Place M[--SP] in D and store it in M[M[13]] (M[VM argument address])
		hack += popToD() + "@R13\nA=M\nM=D"

	} else {
	} //ERROR

	return hack
}

func vmArgumentAddressToAD(command Command) (hack string) {
	switch strings.ToLower(command.arg1) {
	case "local":
		hack = "@LCL\nAD=M\n"
		if arg2 := command.arg2; arg2 == "0" {
		} else if arg2 == "1" {
			hack += "AD=A+1\n"
		} else {
			hack += "@" + arg2 + "\nAD=D+A\n"
		}
	case "argument":
		hack = "@ARG\nAD=M\n"
		if arg2 := command.arg2; arg2 == "0" {
		} else if arg2 == "1" {
			hack += "AD=A+1\n"
		} else {
			hack += "@" + arg2 + "\nAD=D+A\n"
		}
	case "this":
		hack = "@THIS\nAD=M\n"
		if arg2 := command.arg2; arg2 == "0" {
		} else if arg2 == "1" {
			hack += "AD=A+1\n"
		} else {
			hack += "@" + arg2 + "\nAD=D+A\n"
		}
	case "that":
		hack = "@THAT\nAD=M\n"
		if arg2 := command.arg2; arg2 == "0" {
		} else if arg2 == "1" {
			hack += "AD=A+1\n"
		} else {
			hack += "@" + arg2 + "\nAD=D+A\n"
		}
	}
	return hack
}

func arithmeticToHack(command Command) (hack string) {

}

// popToD Place M[--SP] in D
func popToD() string {
	return "@SP\nAM=M-1\nD=M\n"
}

// pushFromD Place D in M[SP++]
func pushFromD() string {
	return "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
}
