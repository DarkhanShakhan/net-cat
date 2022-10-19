package internal

import (
	"bufio"
	"net"
)

type User struct {
	name   string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewUser(conn net.Conn) *User {
	return &User{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}
