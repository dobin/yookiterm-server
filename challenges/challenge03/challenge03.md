# Shellcode Intro Lab

## Introduction

The shellcode is a small piece of code used as the payload in the exploitation of a software vulnerability. It is called "shellcode" because it typically starts a command shell from which the attacker can control the compromised machine, but any piece of code that performs a similar task can be called shellcode. In this challenge we will create a simple shellcode and test it.


## Goal

- Learn more about shellcode creation.
- How to create a shellcode from assembler
- How to test a shellcode


## Preparation

We have a simple assembler program, which should print the message "Hi there" on the console:

```
$ cat print.asm
section .data
msg db 'Hi there',0xa

section .text
global _start
_start:

; write (int fd, char *msg, unsigned int len);
mov eax, 4
mov ebx, 1
mov ecx, msg
mov edx, 9
int 0x80

; exit (int ret)
mov eax, 1
mov ebx, 0
int 0x80
```

Compile it:
```
$ nasm -f elf print.asm  
```

This should generate an object ELF file with the name `print.o`.

Link it with the linker `ld`:
```
$ ld -m elf_i386 -o print print.o  
```

This will generate an executable file "print". Note that we link it as x32, because the source assembler code is in x32. Alternatively just type `make print`.

Try it:
```
$ ./print
Hi there
$  
```

It looks like our code is working.

## Look at it:

We can decompile the generated ELF binary, to check the assembler source code again. Note that the initial program was written in Intel syntax, but objdump will use AT&T syntax.

```
objdump -M intel -d <program> (this will give Intel syntax)
```

```
# objdump -d print

print: file format elf32-i386


Disassembly of section .text:

08048080 <_start>:
8048080:   b8 04 00 00 00   mov  $0x4,%eax
8048085:   bb 01 00 00 00   mov  $0x1,%ebx
804808a:   b9 a4 90 04 08   mov  $0x80490a4,%ecx
804808f:   ba 09 00 00 00   mov  $0x9,%edx
8048094:   cd 80            int  $0x80
8048096:   b8 01 00 00 00   mov  $0x1,%eax
804809b:   bb 00 00 00 00   mov  $0x0,%ebx
80480a0:   cd 80            int  $0x80
```

## Create shellcode

Extract byte-shellcode out of your executable using objdump output
```
$ objdump -d print | grep "^ " \
 | cut -d$'\t' -f 2 | tr '\n' ' ' | sed -e 's/ *$//' \
 | sed -e 's/ \+/\\x/g' | awk '{print "\\x"$0}'

\xb8\x04\x00\x00\x00\xbb\x01\x00\x00\x00\xb9\xa4\x90\x04\x08\xba\x09\x00\x00\x00\xcd\x80\xb8\x01\x00\x00\x00\xbb\x00\x00\x00\x00\xcd\x80
```

The command line above will extract the byte sequence of the program. Sadly we have a lot of 0 bytes `\x00` in the shellcode. We also have a static reference.

We have to remove the 0 bytes, and the static reference, as 0 bytes are considered as string termination, and should not appear in bytecode.


## Remove 0 Bytes from bytecode

New 'print' assembler code that will not contain any 0-bytes: the trick is; use XOR and AL, BL

Get assembler code from print2.asm
```
section .data
msg db 'Hi there',0xa

section .text
global _start
_start:

xor eax,eax
xor ebx,ebx
xor ecx,ecx
xor edx,edx

mov al, 0x4
mov bl, 0x1
mov ecx, msg
mov dl, 0x8
int 0x80

mov al, 0x1
xor ebx,ebx
int 0x80
```

Compile and link it:
```
$ nasm -f elf print2.asm
$ ld -m elf_i386 -o print2 print2.o
```

or build it via `make print2`.

Run it:
```
$ ./print2
Hi there
```

Seems it's still working. But are the 0 bytes removed? Lets check:
```
# objdump -d print2
print2: file format elf32-i386

Disassembly of section .text:
08048080 : <_start>
8048080: 31 c0           xor %eax,%eax
8048082: 31 db           xor %ebx,%ebx
8048084: 31 c9           xor %ecx,%ecx
8048086: 31 d2           xor %edx,%edx
8048088: b0 04           mov $0x4,%al
804808a: b3 01           mov $0x1,%bl
804808c: b9 9c 90 04 08  mov $0x804909c,%ecx
8048091: b2 08           mov $0x8,%dl
8048093: cd 80           int $0x80
8048095: b0 01           mov $0x1,%al
8048097: 31 db           xor %ebx,%ebx
8048099: cd 80           int $0x80
```

Awesome, no more null bytes! But we still have a problem with the hard-coded address $0x804909c (this causes problems in a real shellcode)
Remove References
We need to remove the reference to the data section. For this, we just push the bytes of the message 'hi there' on the stack, and reference them in the stack.
Create BYTES of message "hi there"

```
$ python -c 'print "hi there"' | hexdump -C -v
00000000 68 69 20 74 68 65 72 65 0a |hi there.|
```

And convert it:
```
little endian: 68 65 72 65 --> 65 72 65 68
little endian: 68 69 20 74 --> 74 20 69 68

0a = CR/LF
```

Create new ASM file with built-in message 'hi there'
```
$ cat print3.asm
section .data

section .text
global _start
_start:

xor eax,eax
xor ebx,ebx
xor ecx,ecx
xor edx,edx

mov al, 0x4
mov bl, 0x1
mov dl, 0x8
push 0x65726568
push 0x74206948
mov ecx, esp
int 0x80

mov al, 0x1
xor ebx,ebx
int 0x80
```

Compile and link it:

```
$ nasm -f elf print3.asm
$ ld -o print3 -m elf_i386 print3.o
```

or `make print3`

Try it:
```
$ ./print3
Hi there
```

Dump your shellcode:
```
$ objdump -d print3 | grep "^ " | cut -d$'\t' -f 2 | tr '\n' ' ' | sed -e 's/ *$//' | sed -e 's/ \+/\\x/g'| awk '{print "\\x"$0}'
\x31\xc0\x31\xdb\x31\xc9\x31\xd2\xb0\x04\xb3\x01\xb2\x08\x68\x68\x65\x72\x65\x68\x48\x69\x20\x74\x89\xe1\xcd\x80\xb0\x01\x31\xdb\xcd\x80
```

We can now use this bytecode sequence, and use it in our shellcode test program.

## Test Shellcode with Loader

You can now try the new shellcode in a shellcode loader program

Get print-shellcodetest.c
```
$ cat shellcodetest.c
#include <stdio.h>
#include <string.h>

char *shellcode = "\x31\xc0\x31\xdb\x31\xc9\x31\xd2\xb0\x04\xb3\x01\xb2\x08\x68\x68\x65\x72\x65\x68\x48\x69\x20\x74\x89\xe1\xcd\x80\xb0\x01\x31\xdb\xcd\x80";

int main(void) {
	( *( void(*)() ) shellcode)();
}
$ gcc shellcodetest.c –m32 –z execstack -o shellcodetest
$ ./shellcodetest
Hi there
$
```

## Missions

### Mission 1

Can you make the shellcode smaller? How small?

### Mission 1

Instead of using the system call write(), use the system call 11 (0xb), "sys_execve". Start a bash shell instead of printing 'hi there'

### Mission 2

Do the lab above for 64 bit.


## Questions

- Why are the two move instructions required before triggering the interrupt in print.asm?
