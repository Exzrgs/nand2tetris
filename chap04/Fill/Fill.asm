(LOOP)
    // 0(何も押されていない)だったらWHITEにいく。
    // そうでなければこのままBLACKにいく
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
            // スクリーンのアドレスを格納
            D=A
            @i
            // @valueはAレジスタに格納しているだけ。だから、A=D+Mをしても、Aの値が変わるだけ。iのアドレスは変わらない
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