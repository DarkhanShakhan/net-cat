package internal

import "net"

type Client struct {
	address string
	name    string
	conn    net.Conn
}

var (
	leaving  = make(chan message)
	messages = make(chan message)
)

type message struct {
	text string
	from string
}
