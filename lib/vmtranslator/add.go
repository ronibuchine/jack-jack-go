package vmtranslator

func translateAdd(cmd Command) string {
    /*
    // ADD

    @SP      # Set A = 0
    A=M     # Set A = RAM[0]
    D=M     # Set D = RAM[RAM[0]]
    A=A-1   # Set A = A - 1 (and now M = RAM[A-1])
    M=M+D   # RAM[RAM[0] - 1] += RAM[RAM[0]], or increase second value in the
                #stack by the value of the first position in the stack
    @SP      # Set A=0 and M=RAM[0]
    M=M-1   # Decrement SP


    pop from SP
    add numbers
    push to SP
    */
    hack := "@SP\n"



    return hack
}

