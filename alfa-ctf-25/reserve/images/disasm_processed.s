0000000000000000 <tc_encrypt>:; (struct __sk_buff *skb)
       0:    r6 = r1    ; skb
       1:    r7 = 0x0

       2:    r3 = r10         ; to
       3:    r3 += -0xe       ; to
       4:    r2 = 0x0         ; offset
       5:    r4 = 0xe         ; len
       6:    call 0x1a        ; call bpf_skb_load_bytes(skb, offset=0x0, to=&var_e, len=0xe)
       7:    if r7 s> r0 goto +0x21 <LBB0_7>    ; if r0 < 0 then retrun -1

       8:    r1 = *(u8 *)(r10 - 0x1)
       9:    r1 <<= 0x8
      10:    r2 = *(u8 *)(r10 - 0x2)
      11:    r1 |= r2
      12:    if r1 != 0x8 goto +0x1c <LBB0_7>   ; if r0 < 0 then retrun -1

      13:    r3 = r10         ; to
      14:    r3 += -0x24      ; to
      15:    r1 = r6          ; skb
      16:    r2 = 0xe         ; offset
      17:    r4 = 0x14        ; len
      18:    call 0x1a        ; call bpf_skb_load_bytes(skb, offset=0xe, to=&var_24, len=0x14)
      19:    r1 = 0x0
      20:    if r1 s> r0 goto +0x14 <LBB0_7>    ; if r0 < 0 then retrun -1

      21:    r1 = *(u8 *)(r10 - 0x1b)
      22:    if r1 != 0x6 goto +0x12 <LBB0_7>   ; if r0 < 0 then retrun -1

      23:    r7 = *(u8 *)(r10 - 0x24)
      24:    r7 <<= 0x2
      25:    r7 &= 0x3c
      26:    r9 = r7
      27:    r9 += 0xe
      28:    r3 = r10         ; to
      29:    r3 += -0x38      ; to
      30:    r1 = r6          ; skb
      31:    r2 = r9          ; offset
      32:    r4 = 0x14        ; len
      33:    call 0x1a        ; call bpf_skb_load_bytes(skb, offset, to=&var_38, len=0x14)
      34:    r1 = 0x0
      35:    if r1 s> r0 goto +0x5 <LBB0_7>     ; if r0 < 0 then retrun -1

      36:    r8 = *(u16 *)(r10 - 0x2c)
      37:    r1 = *(u16 *)(r10 - 0x38)
      38:    if r1 == 0x3905 goto +0x5 <LBB0_8>       ; if r1 == 1337 then .main_work
      39:    r1 = *(u16 *)(r10 - 0x36)
      40:    if r1 == 0x3905 goto +0x3 <LBB0_8>       ; if r1 == 1337 then .main_work

0000000000000148 <LBB0_7>:    ; return -1
      41:    r0 = 0xffffffff ll

0000000000000158 <LBB0_18>:   ; retrun r0
      43:    exit

0000000000000160 <LBB0_8>:          ; .main_work
      44:    *(u64 *)(r10 - 0x50) = r9
      45:    r9 = *(u16 *)(r10 - 0x22)

      46:    r2 = *(u32 *)(r6 + 0x0)      ; skb->len
      47:    r2 += 0x4        ; len
      48:    r1 = r6          ; skb
      49:    r3 = 0x0         ; flags
      50:    call 0x26        ; call bpf_skb_change_tail(skb, len=(skb->len+4), flags=0)
      51:    r1 = r0
      52:    r0 = 0x2
      53:    r2 = 0x0
      54:    if r2 s> r1 goto -0xc <LBB0_18>    ; return 2

      55:    r8 >>= 0x2
      56:    r8 &= 0x3c
      57:    r1 = *(u64 *)(r10 - 0x50)
      58:    r1 += r8
      59:    *(u64 *)(r10 - 0x50) = r1
      60:    r8 += r7
      61:    r9 = be16 r9
      62:    r9 -= r8

      63:    call 0x7         ; call bpf_get_prandom_u32(void)
      64:    r8 = r0
      65:    *(u32 *)(r10 - 0x38) = r8    ; var_38 = <r0>

      66:    *(u64 *)(r10 - 0x58) = r9
      67:    r9 += 0x3
      68:    r1 = 0x4
      69:    r2 = r9
      70:    if r1 > r9 goto +0x2a <LBB0_16>
      71:    r1 = 0xfffffffc ll
      73:    r2 &= r1
      74:    r2 >>= 0x2
      75:    *(u64 *)(r10 - 0x48) = r2    ; var_48 = <r2>
      76:    r3 = 0x1
      77:    r7 = *(u64 *)(r10 - 0x50)

0000000000000270 <LBB0_11>:               ; .loop_body_11
      78:    *(u64 *)(r10 - 0x40) = r3    ; var_40=<r3>

      79:    r3 = r10         ; to
      80:    r3 += -0x24      ; to
      81:    r1 = r6          ; skb
      82:    r2 = r7          ; offset
      83:    r4 = 0x4         ; len
      84:    call 0x1a        ; call bpf_skb_load_bytes(skb, offset=<r7>, to=&var_24, len=4)
      85:    r1 = 0x0
      86:    if r1 s> r0 goto +0x29 <LBB0_15>   ; if r0 < 0 then return 2

      ; var_24 ^= <r8>  ; <r8>=current key
      87:    r9 = *(u32 *)(r10 - 0x24)
      88:    r9 ^= r8
      89:    *(u32 *)(r10 - 0x24) = r9

      90:    r3 = r10         ; from
      91:    r3 += -0x24      ; from
      92:    r1 = r6          ; skb
      93:    r2 = r7          ; offset
      94:    r4 = 0x4         ; len
      95:    r5 = 0x0         ; flags
      96:    call 0x9         ; call bpf_skb_store_bytes(skb, offset=<r7>, from=&var_24, len=4, flags=0)
      97:    r1 = 0x0
      98:    if r1 s> r0 goto +0x1d <LBB0_15>   ; if r0 < 0 then return 2

      ; var_40 &= (1 << 32 - 1)
      99:    r3 = *(u64 *)(r10 - 0x40)
     100:    r1 = r3
     101:    r1 <<= 0x20
     102:    r1 >>= 0x20

     103:    r2 = *(u64 *)(r10 - 0x48) ; var_48
     104:    if r1 >= r2 goto +0x8 <LBB0_16>    ; if var_40 >= var_48 then .footer

     105:    r8 *= -0x7789843f         ; <r8> *= 0x88767bc1 ; <r8> = current key * 0x88767bc1
     106:    r9 ^= r8                  ; <r9> ^= r8         ; <r9> = <r8> ^ crypted
     107:    r3 += 0x1
     108:    r7 += 0x4
     109:    r9 *= 0x20763841          ; <r9> *= 0x20763841 ; <r9> *= 0x20763841
     110:    r8 = r9                   ; <r8> = <r9>        ; current key = ((current key * 0x88767bc1) ^ crypted) *= 0x20763841

     111:    r2 = 0x7d0
     112:    if r2 > r1 goto -0x23 <LBB0_11>    ; if r2 > r1 then .loop_body_11

0000000000000388 <LBB0_16>:         ; .footer
     113:    r1 = *(u64 *)(r10 - 0x50)
     114:    r2 = *(u64 *)(r10 - 0x58)
     115:    r2 += r1         ; offset
     116:    r3 = r10         ; from
     117:    r3 += -0x38      ; from
     118:    r7 = 0x0
     119:    r1 = r6          ; skb
     120:    r4 = 0x4         ; len
     121:    r5 = 0x0         ; flags
     122:    call 0x9         ; call bpf_skb_store_bytes(skb, offset=(var_50+var_58), from=&var_38, len=4, flags=0)
     123:    r1 = r0
     124:    r0 = 0x2
     125:    if r7 s> r1 goto -0x53 <LBB0_18>   ; if r0 < 0 then return 2

     126:    r0 = -0x1
     127:    goto -0x55 <LBB0_18>               ; return -1

0000000000000400 <LBB0_15>:
     128:    r0 = 0x2
     129:    goto -0x57 <LBB0_18>               ; return 2
