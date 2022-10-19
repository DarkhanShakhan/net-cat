package main

import (
	"log"
	"net"
	"net-cat/true_client/internal"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
)

var CONN_PORT = "8989"

func main() {
	if len(os.Args) == 2 {
		// FIXME: check errors(nbr of args, correct conn port)
		CONN_PORT = os.Args[1]
	}
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}

	internal.InitGui(conn)
}
