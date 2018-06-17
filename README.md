# lstf

[license]: https://github.com/yuuki/lstf/blob/master/LICENSE

lstf prints `host flows` (aggregated network connection flows to the same source or destination ports) by Linux netlink and enables you to simply grasp the network relationship between localhost and other hosts.

friend: [yuuki/lsconntrack](https://github.com/yuuki/lsconntrack)

## Features

- Distinction of `active open` and `passive open`
- Print also the number of connections of each flows (the absolute values are meaningless)
- Go portability
- JSON support
- TCP support only

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

### JSON format

```shell
$ lstf -n --json | jq -r -M '.'
[
  {
    "direction": "active",
    "local": {
      "Addr": "10.0.1.9",
      "Port": "many"
    },
    "peer": {
      "addr": "10.0.100.1",
      "port": "3306"
    },
    "connections": 20
  },
  {
    "direction": "passive",
    "local": {
      "addr": "10.0.1.9",
      "port": "80"
    },
    "peer": {
      "addr": "10.0.200.1",
      "port": "many"
    },
    "connections": 27
  },
  ...
]
```

## License

[MIT][license]

## Author

[yuuki](https://github.com/yuuki)
