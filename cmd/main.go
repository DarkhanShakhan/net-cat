package main

import (
	"fmt"
	"log"
	"net"
	"net-cat/internal/lobby"
	"os"
	"strconv"
)

const (
	CONN_TYPE = "tcp"
	USAGE     = "[USAGE]: ./TCPChat $port"
)

var (
	CONN_PORT = "8989"
	LOGO      = ""
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println(USAGE)
		return
	} else if len(os.Args) == 2 {
		if !validPort(os.Args[1]) {
			log.Fatal("invalid port")
		}
		CONN_PORT = os.Args[1]
	}

	listener, err := net.Listen(CONN_TYPE, "localhost:"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port %s\n", CONN_PORT)
	lobby := lobby.NewLobby()
	go lobby.ParseSignal()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go lobby.HandleUser(conn)
	}
}

func validPort(port string) bool {
	_, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return len(port) == 4
}
