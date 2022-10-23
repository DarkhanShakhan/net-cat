package service

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TIME_FORMAT   = "2006-01-02 15:04:05"
	LOGO_FILENAME = "cmd/logo.txt"
	LOGO_HOLDER   = "Welcome to TCP Chat!"
)

func GetPrefix(name string) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), name)
}

func ParseLogo() string {
	data, err := os.ReadFile(LOGO_FILENAME)
	if err != nil {
		return LOGO_HOLDER
	}
	return string(data)
}

func ValidInput(input string) bool {
	for _, ch := range input {
		if ch < 32 {
			return false
		}
	}
	res := strings.TrimSpace(input)
	return res != ""
}
