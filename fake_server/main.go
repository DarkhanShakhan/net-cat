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
		h := &handler{conn: conn, writer: bufio.NewWriter(conn), reader: bufio.NewReader(conn)}
		go h.HandleClient()
	}
}

type handler struct {
	conn   net.Conn
	name   string
	writer *bufio.Writer
	reader *bufio.Reader
}

func (h *handler) HandleClient() {
	h.AskName()
	scanner := bufio.NewScanner(h.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "/list" {
			h.writer.WriteString("3 chat(s) available\n")
			h.writer.WriteString("chat1\nchat2\nchat3\n")
			h.writer.Flush()
		} else if msg == "/users" {
			h.writer.WriteString("2 users(s) are online\n")
			h.writer.WriteString("user1\n")
			h.writer.WriteString("user2\n")
			h.writer.Flush()
		}
		fmt.Println(msg)
	}
}

func (h *handler) AskName() {
	h.writer.WriteString(service.ParseLogo() + "\n")
	h.writer.WriteString("Enter your name:")
	h.writer.Flush()
	h.name, _ = h.reader.ReadString('\n')
	fmt.Print(h.name)

}
func Write(conn net.Conn) {

}

func Read(conn net.Conn) {

}
