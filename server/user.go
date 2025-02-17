package main

import "net"

type user struct {
	address net.Conn
	nick    string
}
