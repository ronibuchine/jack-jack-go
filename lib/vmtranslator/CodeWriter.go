package vmtranslator

import (
	"strings"
)

func pushPopToHack(command Command) (hack string) {
	if command.cmdType == C_PUSH {

		// Place VM argument data in D register
		if strings.ToLower(command.arg1) == "constant" {
			hack = "@" + command.arg2 + "\nD=A\n"
		} else {
			hack = vmArgumentAddressToAD(command) + "D=M\n"
		}

		// Place D in M[SP++]
		hack += pushFromD()

	} else if command.cmdType == C_POP {
		// Place VM argument address in M[13] (temp)
		hack += vmArgumentAddressToAD(command) + tempSaveD("R13")

		// Place M[--SP] in D and store it in M[M[13]] (M[VM argument address])
		hack += popToD() + tempToA("R13") + "\nM=D"

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
	hack = popToD()

	switch strings.ToLower(command.arg1) {
	case "add":
		hack += tempSaveD("R13") + popToD()
		hack += tempToA("R13") + "D=D+A\n" + pushFromD()
	case "sub":
		hack += tempSaveD("R13") + popToD()
		hack += tempToA("R13") + "D=D-A\n" + pushFromD()
	case "neg":
		hack += "D=-D\n" + pushFromD()
	case "eq":
		hack += "@SP\nA=M-1\n"  // A points to top of stack (without moving SP)
		hack += "M=M-D\nM=-M\n" // Subtract D from top of stack and invert its boolean value
	case "gt":
	case "lt":
	case "and":
	case "or":
	case "not":
	}

	return hack
}

// popToD Place M[--SP] in D
func popToD() string {
	return "@SP\nAM=M-1\nD=M\n"
}

// pushFromD Place D in M[SP++]
func pushFromD() string {
	return "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
}

func tempSaveD(register string) string {
	return "@" + register + "\nM=D\n"
}
func tempToA(register string) string {
	return "@" + register + "\nA=M\n"
}
