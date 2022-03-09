package vmtranslator

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// popToD Place M[--SP] in D
const popToD = "@SP\nAM=M-1\nD=M\n"

// pushFromD Place D in M[SP++]
const pushFromD = "@SP\nM=M+1\nA=M-1\nM=D\n"

const infiniteLoop = "@ENDLOOP\n(ENDLOOP)\n0;JMP"

func tempSaveD(register string) string {
	return "@" + register + "\nM=D\n"
}
func tempToA(register string) string {
	return "@" + register + "\nA=M\n"
}

// TranslateCommand Returns the hack string of code for a given command struct
func TranslateCommand(cmd *Command) (string, error) {
	switch cmdType := cmd.cmdType; cmdType {
	case CPush, CPop:
		return pushPopToHack(cmd), nil
	case CArithmetic:
		return arithmeticToHack(cmd), nil
	default:
		//ERROR
		return "", errors.New("ERROR, command that was given:\n" + cmd.ToString() + "\nThis command did not translate correctly.")
	}
}

// Given a command struct, this will return a command in hack for a push or a pop
func pushPopToHack(command *Command) (hack string) {
	if command.cmdType == CPush {

		// Place VM argument data in D register
		if strings.ToLower(command.arg1) == "constant" {
			hack = "@" + command.arg2 + "\nD=A\n"
		} else {
			hack = vmArgumentAddressToAD(command) + "D=M\n"
		}

		// Place D in M[SP++]
		hack += pushFromD

	} else if command.cmdType == CPop {
		// Place VM argument address in M[13] (temp)
		hack += vmArgumentAddressToAD(command) + tempSaveD("R13")

		// Place M[--SP] in D and store it in M[M[13]] (M[VM argument address])
		hack += popToD + tempToA("R13") + "\nM=D"

	} else {
		log.Fatal("ERROR, command that was given:\n" + command.ToString() + "\nThis command does not contain valid members to perform a push/pop.")
	}

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
		return "@" + command.arg2 + "\nD=A\n"
	}

	if arg2 := command.arg2; arg2 == "0" {
	} else if arg2 == "1" {
		hack += "AD=A+1\n"
	} else {
		hack += "@" + arg2 + "\nAD=D+A\n"
	}

	return hack
}

var jmpLabel = 0

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
		hack += "D=M-D\n@JMP" + strconv.Itoa(jmpLabel) + "\n"
		if arithmeticType == "eq" {
			hack += "D;JEQ\n"
		} else if arithmeticType == "gt" {
			hack += "D;JGT\n"
		} else {
			hack += "D;JLT\n"
		}
		hack += "@0\nD=A\n@END" + strconv.Itoa(jmpLabel) + "\n0;JMP\n"   // if false, D=0 and jump to END
		hack += "(JMP" + strconv.Itoa(jmpLabel) + ")\n@0\nD=A-1\n"       // if true, D=-1
		hack += "(END" + strconv.Itoa(jmpLabel) + ")\n@SP\nA=M-1\nM=D\n" // M[SP-1] = D
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
