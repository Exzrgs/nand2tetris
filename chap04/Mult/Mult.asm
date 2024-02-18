@i
M = 0
@R2
M = 0

// 方針としては、R0をR1回足す。

(LOOP)
    @i
    D = M
    @R1
    D = M - D
    @END
    D;JLE

    @R0
    D = M
    @R2
    M = M + D
    @i
    M = M + 1
    @LOOP
    0;JMP

(END)
    @END
    0;JMP