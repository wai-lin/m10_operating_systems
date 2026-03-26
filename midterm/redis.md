# The commands are executed in order

`sudo apt update`

`sudo apt install redis`

`sudo systemctl start redis`

`sudo systemctl status redis`
```
...Active...
```

`redis-cli ping`
```
PONG
```

`ps aux | grep redis`
```
redis        234  0.3  0.1  74728 12424 ?        Ssl  13:36   0:00 /usr/bin/redis-server 127.0.0.1:6379
wailina+     608  0.0  0.0   9284  1972 pts/1    S+   13:38   0:00 grep --color=auto redis
```

`ps -o pid,ppid,cmd -p 234`
```
    PID    PPID CMD
    234       1 /usr/bin/redis-server 127.0.0.1:6379
```
PPID is `1` so it's the `systemd` running the redis. It's because `redis` is running via `systemctl`.

`sudo kill 234` will kill the service.
But `systemd` will revive it immediately. We can check it with `ps aux | grep redis`
```
redis        616  0.9  0.1  74728 13620 ?        Ssl  13:45   0:00 /usr/bin/redis-server 127.0.0.1:6379
wailina+     624  0.0  0.0   9284  1964 pts/1    S+   13:45   0:00 grep --color=auto redis
```
So we can only kill the `redis` with `sudo systemctl stop redis`.
```
$ sudo systemctl stop redis
$ redis-cli ping
Could not connect to Redis at 127.0.0.1:6379: Connection refused
```

Checking memory usage.
```
$ cat /proc/635/status


Name:   redis-server
Umask:  0007
State:  S (sleeping)
Tgid:   635
Ngid:   0
Pid:    635
PPid:   1
...
Kthread:        0
VmPeak:    74728 kB
VmSize:    74728 kB
VmLck:         0 kB
VmPin:         0 kB
VmHWM:     13488 kB
VmRSS:      8760 kB
RssAnon:            3840 kB
RssFile:            4920 kB
RssShmem:              0 kB
VmData:    51476 kB
VmStk:       132 kB
VmExe:      2288 kB
VmLib:     13136 kB
VmPTE:       120 kB
VmSwap:        0 kB
...
...
```

Get libraries used by `redis` with `sudo lsof -p 635 | grep '\.so'`

```
redis-ser 635 redis  mem       REG               0,37           1944 /usr/lib/aarch64-linux-gnu/libgpg-error.so.0.34.0 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           2085 /usr/lib/aarch64-linux-gnu/libzstd.so.1.5.5 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1983 /usr/lib/aarch64-linux-gnu/liblzma.so.5.4.5 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1981 /usr/lib/aarch64-linux-gnu/liblz4.so.1.9.4 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1928 /usr/lib/aarch64-linux-gnu/libgcrypt.so.20.4.3 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1879 /usr/lib/aarch64-linux-gnu/libcap.so.2.66 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1926 /usr/lib/aarch64-linux-gnu/libgcc_s.so.1 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           2053 /usr/lib/aarch64-linux-gnu/libstdc++.so.6.0.33 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1874 /usr/lib/aarch64-linux-gnu/libc.so.6 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1886 /usr/lib/aarch64-linux-gnu/libcrypto.so.3 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           2051 /usr/lib/aarch64-linux-gnu/libssl.so.3 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           2055 /usr/lib/aarch64-linux-gnu/libsystemd.so.0.38.0 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1984 /usr/lib/aarch64-linux-gnu/libm.so.6 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37          78239 /usr/lib/aarch64-linux-gnu/libjemalloc.so.2 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37          78254 /usr/lib/aarch64-linux-gnu/liblzf.so.1.5 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37          26414 /usr/lib/aarch64-linux-gnu/libatomic.so.1.2.0 (path dev=0,40)
redis-ser 635 redis  mem       REG               0,37           1835 /usr/lib/aarch64-linux-gnu/ld-linux-aarch64.so.1 (path dev=0,40)```



