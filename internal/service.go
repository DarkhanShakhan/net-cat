package internal

import (
	"fmt"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

func getPrefix(client Client) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), client.name)
}
