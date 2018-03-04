# lstf

lstf prints `host flows` (aggregated network connection flows to the same source or destination ports) by Linux `/proc/net/tcp` (`netstat -tan`) and enables you to simply grasp the network relationship between localhost and other hosts.

friend: [yuuki/lsconntrack](https://github.com/yuuki/lsconntrack)

## How to use

```shell
$ lstf -n
Local Address:Port   <-->   Peer Address:Port     Connections
localhost:many       -->    10.0.1.10:3306        22
localhost:many       -->    10.0.1.11:3306        14
localhost:many       -->    10.0.1.20:8080        99
localhost:80         <--    10.0.2.10:80          120
localhost:80         <--    10.0.2.11:80          202
```