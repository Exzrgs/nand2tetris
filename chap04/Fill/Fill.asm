(LOOP)
    @KBD
    D=M
    @WHITE
    D;JEQ

    (BLACK)
        @SCREEN
        D=M
        @LOOP
        D;JNE

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

            @8192 //256*32
            D=A
            @i
            D=D-M
            @LOOP2
            D;JNE
        
        @LOOP
        0;JMP

    (WHITE)
        @SCREEN
        D=M
        @LOOP
        D;JEQ

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

            @8192
            D=A
            @i
            D=D-M
            @LOOP3
            D;JNE
        
        @LOOP
        0;JMP