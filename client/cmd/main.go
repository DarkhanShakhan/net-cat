package main

import (
	"log"
	"net"
	"net-cat/client/internal"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
	USAGE     = "go run client/cmd/main.go port"
)

var CONN_PORT = "8989"

func main() {
	if len(os.Args) == 2 {
		// FIXME: check errors(nbr of args, correct conn port)
		CONN_PORT = os.Args[1]
	}
	if len(os.Args) > 2 {
		log.Fatal("too many arguments\n" + USAGE)
	}
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	internal.InitGui(conn)
}
