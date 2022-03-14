package vmtranslator

import (
	"strconv"
	"strings"
)

const bootstrap = "// Bootstrap\n@256\nD=A\n@SP\nM=D\n@5\nD=A\n@SP\nM=M+D\n@Sys.init\n0;JMP\n" // Sets SP to 256 and simulates call to Sys.init but without return value and address

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
func TranslateCommand(cmd *Command) (hack string) {
	switch cmd.cmdType {
	case CPush, CPop:
		hack = pushPopToHack(cmd)
	case CArithmetic:
		hack = arithmeticToHack(cmd)
	case CGoto:
		hack = gotoToHack(cmd)
	case CIfGoto:
		hack = ifGotoToHack(cmd)
	case CLabel:
		hack = labelToHack(cmd)
	case CFunction:
		hack = functionToHack(cmd)
	case CCall:
		hack = callToHack(cmd)
	case CReturn:
		hack = returnToHack()
	}
	return
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
		hack += popToD + tempToA("R13") + "M=D\n"

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
		hack += "D=M-D\n@COMPJUMP" + strconv.Itoa(jmpLabel) + "\n"
		if arithmeticType == "eq" {
			hack += "D;JEQ\n"
		} else if arithmeticType == "gt" {
			hack += "D;JGT\n"
		} else {
			hack += "D;JLT\n"
		}
		hack += "@0\nD=A\n@COMPEND" + strconv.Itoa(jmpLabel) + "\n0;JMP\n"   // if false, D=0 and jump to END
		hack += "(COMPJUMP" + strconv.Itoa(jmpLabel) + ")\n@0\nD=A-1\n"      // if true, D=-1
		hack += "(COMPEND" + strconv.Itoa(jmpLabel) + ")\n@SP\nA=M-1\nM=D\n" // M[SP-1] = D
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

// Takes a vm command of type goto and returns the hack code as a string
func gotoToHack(cmd *Command) (hack string) {

	// store the label in A and jump unconditionally
	hack += "@" + cmd.arg1 + "\n0;JMP\n"
	return

}

// Takes a vm command of type label and returns the hack code as a string
func labelToHack(cmd *Command) (hack string) {

	// creates label of arg1 and wraps in parentheses
	hack += "(" + cmd.arg1 + ")\n"
	return
}

// Takes a vm command of type if-goto and returns the hack code as a string
func ifGotoToHack(cmd *Command) (hack string) {

	// checks the value at the top of the stack and compares it to 0, if yes it jumps to the label
	hack += popToD + "@" + cmd.arg1 + "\nD;JNE\n"

	return

}

func functionToHack(cmd *Command) (hack string) {
	hack = "(" + cmd.arg1 + ")\n" // Inserts function entry label
	if cmd.arg2 == "0" {
		return hack
	}
	hack += "@" + cmd.arg2 + "\nD=A\n@R13\nM=D\n"                                                              // Store number of local variables in R13
	hack += "(" + cmd.arg1 + "$init)\n@0\nD=A\n" + pushFromD + "@R13\nMD=M-1\n@" + cmd.arg1 + "$init\nD;JNE\n" // Push 0 onto stack until R13 is 0
	return hack
}

var funcCallCounter = 0

func callToHack(cmd *Command) (hack string) {
	if cmd.arg2 == "0" && cmd.arg1 != "Sys.init" {
		cmd.arg2 = "1"
		hack = "@SP\nM=M+1\n"
	} // Ensures space for return value
	hack += "@" + cmd.arg1 + "$ret." + strconv.Itoa(funcCallCounter) + "\nD=A\n" + pushFromD // Creates a return label and pushes it onto stack
	hack += "@LCL\nD=M\n" + pushFromD                                                        // Pushes LCL onto stack
	hack += "@ARG\nD=M\n" + pushFromD                                                        // Pushes ARG onto stack
	hack += "@THIS\nD=M\n" + pushFromD                                                       // Pushes THIS onto stack
	hack += "@THAT\nD=M\n" + pushFromD                                                       // Pushes THAT onto stack
	hack += "@SP\nD=M\n@LCL\nM=D\n"                                                          // Sets new LCL segment to top of stack
	hack += "@5\nD=D-A\n" + "@" + cmd.arg2 + "\nD=D-A\n@ARG\nM=D\n"                          // Sets new ARG segment to (top of stack) - 5 - (number of args)
	hack += "@" + cmd.arg1 + "\n0;JMP\n"                                                     // Jump to function
	hack += "(" + cmd.arg1 + "$ret." + strconv.Itoa(funcCallCounter) + ")\n"                 // Label for function to return to caller
	funcCallCounter++                                                                        // Increment call counter

	return hack
}

func returnToHack() (hack string) {
	hack = "@LCL\nD=M\n@R13\nM=D\n"           // Stores top of virtual stack in R13
	hack += popToD + "@ARG\nA=M\nM=D\n"       // Saves top of real stack to ARG[0] (returned value for the caller)
	hack += "D=A+1\n@SP\nM=D\n"               // Sets SP to the value following the returned value (for the caller)
	hack += "@R13\nAM=M-1\nD=M\n@THAT\nM=D\n" // Pushes top of virtual stack to THAT
	hack += "@R13\nAM=M-1\nD=M\n@THIS\nM=D\n" // Pushes top of virtual stack to THIS
	hack += "@R13\nAM=M-1\nD=M\n@ARG\nM=D\n"  // Pushes top of virtual stack to ARG
	hack += "@R13\nAM=M-1\nD=M\n@LCL\nM=D\n"  // Pushes top of virtual stack to LCL
	hack += "@R13\nA=M-1\nA=M\n0;JMP\n"       // Jumps to location at top of stack

	return hack
}
