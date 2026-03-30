wailinaung@Ubuntu:~/m10_operating_systems/hacking$ ps -a
PID TTY TIME CMD
17889 pts/2 00:00:00 zellij
18010 pts/4 00:00:00 password.arm64.
18011 pts/3 00:00:00 ps

wailinaung@Ubuntu:~/m10_operating_systems/hacking$ PID=18010

wailinaung@Ubuntu:~/m10_operating_systems/hacking$ cat /proc/18010/maps | tee password.maps
aaaad22b0000-aaaad22b1000 r-xp 00000000 00:25 78673 /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf
aaaad22c0000-aaaad22c1000 r--p 00000000 00:25 78673 /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf
aaaad22c1000-aaaad22c2000 rw-p 00001000 00:25 78673 /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf
aaaaf9023000-aaaaf9044000 rw-p 00000000 00:00 0 [heap]
ffff83120000-ffff832ba000 r-xp 00000000 00:25 11874 /usr/lib/aarch64-linux-gnu/Libc.so.6
ffff832ba000-ffff832cd000 ---p 0019a000 / 00:25 1874 /usr/lib/aarch64-linux-gnu/libc.so.6
ffff832cd000-ffff832d0000 r--p 0019d000 00:25 1874 /usr/Lib/aarch64-linux-gnu/libc.so.6
ffff832d0000-ffff832d2000 rw-p 001a0000 00:25 1874 /usr/lib/aarch64-linux-gnu/libc.so.6
ffff832d2000-ffff832de000 rw-p : 00000000 00:00 0
ffff832ee000-ffff83315000 r-xp 00000000 | 00:25 1835 /usr/Lib/aarch64-linux-gnu/ld-linux-aarch64.so.1
ffff83324000-ffff83326000 rw-p 00000000 00:00 0
ffff83326000-ffff8332a000 r--p 00000000 00:00 0 [vvar]
ffff8332a000-ffff8332c000 r-xp 00000000 00:00 0 [vdso]
ffff8332c000-ffff8332e000 r--p 0002e000 00:25 1835 /usr/Lib/aarch64-linux-gnu/ld-linux-aarch64.so.1
ffff8332000-ffff83330000 rw-p 00030000 00:25 1835 /usr/Lib/aarch64-linux-gnu/ld-linux-aarch64.so.1
ffffc80a6000-ffffc80c7000 rw-p 00000000 00:00 0 [stack]

wailinaung@Ubuntu:~/m10_operating_systems/hacking$ sudo gdb -q -p $(echo $PID) \
> -ex "set pagination off" \
> -ex "dump memory ./heap.bin 0xaaaaf9023000 0xaaaaf9044000" \
> -ex detach -ex quit

Attaching to process 18010
Reading symbols from /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf...
(No debugging symbols found in /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf)
Reading symbols from /lib/aarch64-linux-gnu/libc.so.6...
Reading symbols from /usr/lib/debug/.build-id/d5/ef86dde36cbd3289566cf5098226035d76f2e1.debug...
Reading symbols from /lib/ld-linux-aarch64.so.1...
Reading symbols from /usr/lib/debug/.build-id/f3/d28c5cab7887a8195f6b130d76b8faf126b168.debug...
[Thread debugging using libthread_db enabled]
Using host libthread_db library "/lib/aarch64-linux-gnu/libthread_db.so.1".
0x0000ffff83201244 in __GI___libc_read (fd=0, buf=0xaaaaf90232c0, nbytes=1024) at ../sysdeps/unix/sysv/linux/read.c:26
warning: 26	../sysdeps/unix/sysv/linux/read.c: No such file or directory
Detaching from program: /home/wailinaung/m10_operating_systems/hacking/password.arm64.elf, process 18010
[Inferior 1 (process 18010) detached]

wailinaung@Ubuntu:~/m10_operating_systems/hacking$ cat heap.bin 
pwAW96B6

wailinaung@Ubuntu:~/m10_operating_systems/hacking$ printf 'pwAW96B6\n' | ./password.arm64.elf 
