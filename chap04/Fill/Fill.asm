(LOOP)
    @KBD
    D=M
    @WHITE
    D;JEQ

    (BLACK)
        @i
        M=0
        (LOOP2)
            @SCREEN
            D=A
            @i
            A=D+M
            M=-1
            
            @i
            M=M+1

            @8000
            D=A
            @i
            D=D-M
            @LOOP2
            D;JGT
        
        @LOOP
        0;JMP

    (WHITE)
        @i
        M=0
        (LOOP3)
            @SCREEN
            D=A
            @i
            A=D+M
            M=0
            
            @i
            M=M+1

            @8000
            D=A
            @i
            D=D-M
            @LOOP3
            D;JGT
        @LOOP
        0;JMP