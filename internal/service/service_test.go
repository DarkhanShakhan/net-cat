package service

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPrefix(t *testing.T) {
	for i := 1; i <= 100; i++ {
		expPrefix := fmt.Sprintf("[%s][user_%d]:", time.Now().Format(TIME_FORMAT), 1)
		resPrefix := GetPrefix(fmt.Sprintf("user_%d", 1))
		if expPrefix != resPrefix {
			t.Errorf("'GetPrefix' function FAILED, expected %s, got %s", expPrefix, resPrefix)
		} else {
			t.Logf("'GetPrefix' function SUCCEEDED, expected %s, got %s", expPrefix, resPrefix)
		}
	}
}
