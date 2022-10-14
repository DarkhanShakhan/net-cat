package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	CONN_PORT = ":8989"
	CONN_TYPE = "tcp"

	MSG_DISCONNECT = "Disconnected from the server.\n"
)

var wg sync.WaitGroup

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Bytes()

		// str, err := reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Printf(MSG_DISCONNECT)
		// 	wg.Done()
		// 	return
		// }
		fmt.Print(string(text))
	}
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = writer.WriteString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	// wg.Add(1)
	conn, err := net.Dial(CONN_TYPE, CONN_PORT)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	PrintLogo(conn)
	// PrintName(conn)
	// go Read(conn)
	// go Write(conn)
	// wg.Wait()
}

func PrintLogo(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for i := 0; i < 17; i++ {
		line, _ := reader.ReadString('\n')
		fmt.Print(line)
	}
	askname, _ := reader.ReadString(':')
	fmt.Print(askname)
	readstdin := bufio.NewReader(os.Stdin)
	name, _ := readstdin.ReadString('\n')
	conn.Write([]byte(name))
}

// func PrintName(conn net.Conn) {
// 	reader := bufio.NewReader(conn)
// 	readstdin := bufio.NewReader(os.Stdin)
// 	askname, err := reader.ReadString('\n')
// 	fmt.Println(err)
// 	fmt.Println("here")
// 	fmt.Print(askname)
// 	name, _ := readstdin.ReadString('\n')
// 	conn.Write([]byte(name))
// }
