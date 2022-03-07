package vmtranslator

import (
	"log"
	"strconv"
	"strings"
)

func TranslateCommand(cmd *Command) (string, error) {
	log.Fatal("unimplemented")
	switch cmdType := cmd.cmdType; cmdType {
	case C_PUSH, C_POP:
		return pushPopToHack(cmd), nil
	case C_ARITHMETIC:
		return arithmeticToHack(cmd), nil
	default:
		//ERROR
		return "", nil
	}
}

func pushPopToHack(command *Command) (hack string) {
	if command.cmdType == C_PUSH {

		// Place VM argument data in D register
		if strings.ToLower(command.arg1) == "constant" {
			hack = "@" + command.arg2 + "\nD=A\n"
		} else {
			hack = vmArgumentAddressToAD(command) + "D=M\n"
		}

		// Place D in M[SP++]
		hack += pushFromD

	} else if command.cmdType == C_POP {
		// Place VM argument address in M[13] (temp)
		hack += vmArgumentAddressToAD(command) + tempSaveD("R13")

		// Place M[--SP] in D and store it in M[M[13]] (M[VM argument address])
		hack += popToD + tempToA("R13") + "\nM=D"

	} else {
	} //ERROR

	return hack
}

func vmArgumentAddressToAD(command *Command) (hack string) {
	switch command.arg1 {
	case "local":
		hack = "@LCL\nAD=M\n"
	case "argument":
		hack = "@ARG\nAD=M\n"
	case "this":
		hack = "@THIS\nAD=M\n"
	case "that":
		hack = "@THAT\nAD=M\n"
	case "temp":
		hack = "@5\nD=A\n"
	case "pointer":
		hack = "@THIS\nD=A\n"
	case "static":
		return "@" + command.arg2 + "\nD=A\n" //TODO fix filename situation
	}

	if arg2 := command.arg2; arg2 == "0" {
	} else if arg2 == "1" {
		hack += "AD=A+1\n"
	} else {
		hack += "@" + arg2 + "\nAD=D+A\n"
	}

	return hack
}

var jmpLabel int = 0

func arithmeticToHack(command *Command) (hack string) {

	switch arithmeticType := strings.ToLower(command.arg1); arithmeticType {
	case "add":
		hack = popToD + "@SP\nA=M-1\nM=M+D\n"
	case "sub":
		hack = popToD + "@SP\nA=M-1\nM=M-D\n"
	case "neg":
		hack = "@SP\nA=M-1\nM=-M\n"
	case "eq", "gt", "lt":
		hack = popToD
		hack += "@SP\nA=M-1\n" // A points to top of stack (without moving SP)
		hack += "M=M-D\n@-1\nD=A\n@JMP" + strconv.Itoa(jmpLabel) + "\n"
		if arithmeticType == "eq" {
			hack += "M;JEQ"
		} else if arithmeticType == "gt" {
			hack += "M;JGT\n"
		} else {
			hack += "M;JLT\n"
		}
		hack += "@0\nD=A\n(JMP" + strconv.Itoa(jmpLabel) + ")\n" + "@SP\nA=M-1\nM=D\n"
		jmpLabel += 1
	case "and":
		hack = popToD
		hack += "@SP\nA=M-1\nM=D&M\n"
	case "or":
		hack = popToD
		hack += "@SP\nA=M-1\nM=D|M\n"
	case "not":
		hack = "@SP\nA=M-1\nM=!M\n"
	}

	return hack
}

// popToD Place M[--SP] in D
const popToD = "@SP\nAM=M-1\nD=M\n"

// pushFromD Place D in M[SP++]
const pushFromD = "@SP\nM=M+1\nA=M-1\nM=D\n"

func tempSaveD(register string) string {
	return "@" + register + "\nM=D\n"
}
func tempToA(register string) string {
	return "@" + register + "\nA=M\n"
}
