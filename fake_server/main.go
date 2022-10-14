package main

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/internal/service"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		h := &handler{conn: conn}
		go h.HandleClient()
	}
}

type handler struct {
	conn net.Conn
	name string
}

func (h *handler) HandleClient() {
	h.AskName()
	scanner := bufio.NewScanner(h.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println(msg)
	}
}

func (h *handler) AskName() {
	h.conn.Write([]byte(service.ParseLogo() + "\n"))
	h.conn.Write([]byte("Enter your name:"))
	h.name, _ = bufio.NewReader(h.conn).ReadString('\n')
	h.conn.Write([]byte("Welcome, " + h.name))

}
func Write(conn net.Conn) {

}

func Read(conn net.Conn) {

}
