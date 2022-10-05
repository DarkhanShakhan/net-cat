package internal

import (
	"fmt"
)

func Broadcaster(clients map[string]Client) {
	for {
		select {
		case msg := <-messages:
			for _, client := range clients {
				if msg.from == client.conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(client.conn, "")
				fmt.Fprintln(client.conn, msg.text)
				fmt.Fprintf(client.conn, getPrefix(client))
			}
		case msg := <-leaving:
			for _, client := range clients {
				fmt.Fprintln(client.conn, "")
				fmt.Fprintln(client.conn, msg.text)
				fmt.Fprintf(client.conn, getPrefix(client))
			}
		}
	}
}
