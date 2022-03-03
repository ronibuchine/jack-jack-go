package vmtranslator

func translateAdd() string {
	/*
	   // ADD

	   @SP      # Set A = 0
	   A=M      # Set A = RAM[0]
	   D=M      # Set D = RAM[RAM[0]]
	   A=A-1    # Set A = A - 1 (and now M = RAM[A-1])
	   M=M+D    # RAM[RAM[0] - 1] += RAM[RAM[0]], or increase second value in the
	               #stack by the value of the first position in the stack
	   @SP      # Set A=0 and M=RAM[0]
	   M=M-1    # Decrement SP
	*/
	return "@SP\nA=M\nD=M\nA=A-1\nM=M+D\n@SP\nM+M-1\n"
}
