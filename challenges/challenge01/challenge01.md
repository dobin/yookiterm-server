# Buffer Overflow Analysis Intro Lab

This little practice lab is geared toward beginners and those who want to get introduced into the eco-system of a penetration tester or ethical hacker.

You will not exploit a buffer overflow in this lab; instead you get introduced to some basic tools and knowledge. This lab shall be done before the buffer overflow labs.

## source

```
#include <stdio.h>
#include <stdlib.h>

char globalVariable[] = "GlobalVar";

int main(int argc, char **argv) {
    if (argc == 1) {
        printf("Call: %s <name>\n", argv[0]);
        exit(0);
    }

    printf("Hello %s\n", argv[1]);
}
```


## Check File Types of the new binaries

Next, let's check if you really got a 32bit and 64bit binary out of the two gcc calls above

```
$ file 7377_bof 7377m32_bof
7377_bof:    ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 2.6.32, BuildID[sha1]=ea748ed128f6cc70a0e496f4c3592a32e4323404, not stripped
7377m32_bof: ELF 32-bit LSB executable, Intel 80386, version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux.so.2, for GNU/Linux 2.6.32, BuildID[sha1]=ca4c0d02ce9809aada956abb9db6ccb980ed4a7b, not stripped
```



Run the binaries and test short and long arguments on the command line

Play around with the arguments.

```
./7377_bof hello

./7377m32_bof cool


Instead of manually giving the argument, python can do the job for you

─$ ./7377_bof `python -c 'print "A"*100'`
Hello AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA

./7377_bof $(perl -e 'print "A"x100')
Hello AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
```


Analyze the binary using objdump

Try objdump
```

objdump -d 7377_bof
```


Type readelf
```

readelf -l 7377_bof
readelf -d 7377_bof
readelf -d 7377_bof
Dynamic section at offset 0x788 contains 24 entries:   
Tag        Type                         Name/Value  
0x0000000000000001 (NEEDED)             Shared library: [libc.so.6]
0x000000000000000c (INIT)               0x4003e0  
0x000000000000000d (FINI)               0x400614  
0x0000000000000019 (INIT_ARRAY)         0x600770  
0x000000000000001b (INIT_ARRAYSZ)       8 (bytes)  
0x000000000000001a (FINI_ARRAY)         0x600778  
0x000000000000001c (FINI_ARRAYSZ)       8 (bytes)  
0x000000006ffffef5 (GNU_HASH)           0x400260  
0x0000000000000005 (STRTAB)             0x4002f8  
0x0000000000000006 (SYMTAB)             0x400280  
0x000000000000000a (STRSZ)              68 (bytes)  
0x000000000000000b (SYMENT)             24 (bytes)  
0x0000000000000015 (DEBUG)              0x0  
0x0000000000000003 (PLTGOT)             0x600960  
0x0000000000000002 (PLTRELSZ)           96 (bytes)  
0x0000000000000014 (PLTREL)             RELA  
0x0000000000000017 (JMPREL)             0x400380  
0x0000000000000007 (RELA)               0x400368  
0x0000000000000008 (RELASZ)             24 (bytes)  
0x0000000000000009 (RELAENT)            24 (bytes)  
0x000000006ffffffe (VERNEED)            0x400348  
0x000000006fffffff (VERNEEDNUM)         1  
0x000000006ffffff0 (VERSYM)             0x40033c  
0x0000000000000000 (NULL)               0x0
```


GDB Info Functions

Let's debug the binary using gdb; listing functions

```

gdb ./7377m32_bof

(gdb) info func (lists all functions in the binary)
```





PLT stands for Procedure Linkage Table which is, put simply, used to call external procedures/functions whose address isn't known in the time of linking, and is left to be resolved by the dynamic linker at run time.



GDB Disassemble main in binary


Now let's disassemble main

```
(gdb) disass main

GDB Breakpoint *main and get 20hex values of $ESP info

With this last step, we want to run the binary in gdb


gdb ./7377m32_bof
(gdb) break *main
(gdb) run
(gdb) x/20x $esp
(gdb) c
```





Turn ASLR on

If you have finished this task, please turn ASLR on on your computer

echo 1 > /proc/sys/kernel/randomize_va_space
Security Questions
Please respond to the following security questions

What means plt?
Explain how you create a breakpoint in gdb at the main routine
Explain how you disclose the ESP at a given breakpoint
Explain how to continue debugging after you have seen the details of a breakpoint
