package service

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	TIME_FORMAT   = "2006-01-02 15:04:05"
	LOGO_FILENAME = "cmd/logo.txt"
)

func GetPrefix(name string) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), name)
}

func ParseLogo() string {
	data, err := os.ReadFile(LOGO_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
