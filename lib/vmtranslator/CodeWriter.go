package vmtranslator

import (
	"os"
	"strings"
)

func writePushPop(command Command_t, file os.File) {
	var hack string

	if segment := command._cmdType; segment == C_PUSH {

		// Place VM argument data in D register
		switch strings.ToLower(command._arg1) {
		case "local":
			hack = "@LCL\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "A=A+1\nD=M\n"
			} else {
				hack += "@" + arg2 + "\nA=D+A\nD=M\n"
			}
		case "argument":
			hack = "@ARG\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "A=A+1\nD=M\n"
			} else {
				hack += "@" + arg2 + "\nA=D+A\nD=M\n"
			}
		case "this":
			hack = "@THIS\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "A=A+1\nD=M\n"
			} else {
				hack += "@" + arg2 + "\nA=D+A\nD=M\n"
			}
		case "that":
			hack = "@THAT\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "A=A+1\nD=M\n"
			} else {
				hack += "@" + arg2 + "\nA=D+A\nD=M\n"
			}
		case "constant":
			hack = "@" + command._arg2 + "\nA=M\n"
		default: //ERROR
		}

		// Place D in M[SP]
		hack += "@SP\nA=M\nM=D\n"

		// Increment SP
		hack += "@SP\nM=M+1\n"

	} else if segment == C_POP {
		var hack string

		// Place VM argument address in D register
		switch strings.ToLower(command._arg1) {
		case "local":
			hack = "@LCL\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "D=D+1\n"
			} else {
				hack += "@" + arg2 + "\nD=D+A\n"
			}
		case "argument":
			hack = "@ARG\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "D=D+1\n"
			} else {
				hack += "@" + arg2 + "\nD=D+A\n"
			}
		case "this":
			hack = "@THIS\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "D=D+1\n"
			} else {
				hack += "@" + arg2 + "\nD=D+A\n"
			}
		case "that":
			hack = "@THAT\nD=M\n"
			if arg2 := command._arg2; arg2 == '0' {
			} else if arg2 == '1' {
				hack += "D=D+1\n"
			} else {
				hack += "@" + arg2 + "\nD=D+A\n"
			}
		default: //ERROR
		}

		// Decrement SP and place M[SP] in M[D]
		hack += "@SP\nAM=M-1\nM=D\n"
	} else {
	} //ERROR
	file.WriteString(hack)
}
