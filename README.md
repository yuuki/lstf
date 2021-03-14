# lstf
[![Latest Version](http://img.shields.io/github/release/yuuki/lstf.svg?style=flat-square)](https://github.com/yuuki/lstf/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuuki/lstf)](https://goreportcard.com/report/github.com/yuuki/lstf)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

lstf prints `host flows` (aggregated network connection flows to the same source or destination ports) by Linux netlink and enables you to simply grasp the network relationship between localhost and other hosts.

friend: [yuuki/lsconntrack](https://github.com/yuuki/lsconntrack)

## Features

- Distinction of `active open` and `passive open`
- Print also the number of connections of each flows (the absolute values are meaningless)
- Go portability
- JSON support
- TCP support only

## Installation

### Download binary from GitHub Releases

<https://github.com/yuuki/lstf/releases>

## How to use

HTTP requests --> Web:80 --> MySQL:3306

```shell
$ lstf -n
Local Address:Port   <-->   Peer Address:Port     Connections
10.0.1.9:many        -->    10.0.1.10:3306        22
10.0.1.9:many        -->    10.0.1.11:3306        14
10.0.2.10:22         <--    192.168.10.10:many    1
10.0.1.9:80          <--    10.0.2.13:many        120
10.0.1.9:80          <--    10.0.2.14:many        202
```

- `-->` indicates `active open`
- `<--` indicates `passive open`

Sort flows by the number of connection.

```shell
$ lstf -n | sort -nrk4
```

### JSON format

```shell-session
$ lstf --json | jq -r -M '.'
[
  {
    "direction": "active",
    "local": {
      "name"| "app01.local",
      "addr": "10.0.1.9",
      "port": "many"
    },
    "peer": {
      "name"| "db01.local",
      "addr": "10.0.100.1",
      "port": "3306"
    },
    "connections": 20
  },
  {
    "direction": "passive",
    "local": {
      "name"| "app01.local",
      "addr": "10.0.1.9",
      "port": "80"
    },
    "peer": {
      "name"| "web01.local",
      "addr": "10.0.200.1",
      "port": "many"
    },
    "connections": 27
  },
  ...
]
```

## License

[MIT](LICENSE)

## Author

[yuuki](https://github.com/yuuki)
