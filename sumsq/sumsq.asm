;2**2 + 10**2 = 104

START   LDA TWO
        STA COUNT
        LDA ZERO
        STA ACCUM
SQ2LP   LDA COUNT
        BRZ SQ2END
        LDA ACCUM
        ADD TWO
        STA ACCUM
        LDA COUNT
        SUB ONE
        STA COUNT
        BRA SQ2LP
SQ2END  LDA ACCUM
        STA SQ2

        LDA TEN
        STA COUNT
        LDA ZERO
        STA ACCUM
SQ10LP  LDA COUNT
        BRZ SQ10END
        LDA ACCUM
        ADD TEN
        STA ACCUM
        LDA COUNT
        SUB ONE
        STA COUNT
        BRA SQ10LP
SQ10END LDA ACCUM
        STA SQ10

        LDA SQ2
        ADD SQ10
        OUT
        HLT

; Data
ONE     DAT 1
TWO     DAT 2
TEN     DAT 10
ZERO    DAT 0
COUNT   DAT 0
ACCUM   DAT 0
SQ2     DAT 0
SQ10    DAT 0
