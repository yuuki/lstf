package netutil

import (
	"net"
	"strings"
)

// UserEnt represents a detail of network socket.
// see https://github.com/shemminger/iproute2/blob/afa588490b7e87c5adfb05d5163074e20b6ff14a/misc/ss.c#L509.
type UserEnt struct {
	inode uint32 // inode number
	fd    int    // file discryptor
	pid   int    // process id
	pname string // process name
	ppid  int    // parent process id
	pgrp  int    // process group id
}

// Inode returns inode.
func (u *UserEnt) Inode() uint32 {
	return u.inode
}

// Fd returns file descriptor.
func (u *UserEnt) Fd() int {
	return u.fd
}

// Pid returns process id.
func (u *UserEnt) Pid() int {
	return u.pid
}

// Pname returns process name.
func (u *UserEnt) Pname() string {
	return u.pname
}

// Ppid returns process id.
func (u *UserEnt) Ppid() int {
	return u.ppid
}

// Pgrp returns process group id.
func (u *UserEnt) Pgrp() int {
	return u.pgrp
}

func (u *UserEnt) SetInode(inode uint32) {
	u.inode = inode
}

// UserEnts represents a hashmap of UserEnt as key is the inode.
type UserEnts map[uint32]*UserEnt

// ResolveAddr lookup first hostname from IP Address.
func ResolveAddr(addr string) string {
	hostnames, _ := net.LookupAddr(addr)
	if len(hostnames) > 0 {
		return strings.TrimSuffix(hostnames[0], ".")
	}
	return addr
}

// LocalIPAddrs gets the string slice of localhost IPaddrs.
func LocalIPAddrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	addrStrings := make([]string, 0, len(addrs))
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addrStrings = append(addrStrings, ipnet.IP.String())
			}
		}
	}
	return addrStrings, nil
}
