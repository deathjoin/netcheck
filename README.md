# netcheck

Simple script to check if host alive by sending ICMP messages & TCP Port checks.

ICMP messages not working without sudo privileges.

## usage

```
netcheck -help
Usage of netcheck:
  -host string
        Server ip or name to check.
  -interval int
        Check interval in seconds. (default 1)
  -only_errors
        Print only fails.
  -port int
        Server TCP port to check.
  -timeout int
        Connection timeout in seconds. (default 5)
```

## examples

```
 sudo netcheck -host ya.ru -port 80
2021/10/13 16:54:52 Starting tcp port check: ya.ru:80
2021/10/13 16:54:52 Connection success to "ya.ru:80"
2021/10/13 16:54:54 Connection success to "ya.ru:80"
^C
❯ sudo netcheck -host ya.ru -port 3366 -only_errors
2021/10/13 16:55:04 Starting tcp port check: ya.ru:3366
ERROR: 2021/10/13 16:55:09 Connection fail to "ya.ru:3366" within 5 seconds. PING: OK, TCP_PORT: FAIL
ERROR: 2021/10/13 16:55:16 Connection fail to "ya.ru:3366" within 5 seconds. PING: OK, TCP_PORT: FAIL
^C
❯ sudo netcheck -host 192.168.8.119 -port 22
2021/10/13 16:56:56 Starting tcp port check: 192.168.8.119:22
2021/10/13 16:56:56 Connection success to "192.168.8.119:22"
2021/10/13 16:56:57 Connection success to "192.168.8.119:22"
2021/10/13 16:56:58 Connection success to "192.168.8.119:22"
^C
```