package user

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestUser(t *testing.T) {

	for i := 1; i <= 10; i++ {
		connName, userName, roomName := fmt.Sprintf("tcpConn_%d", i), fmt.Sprintf("user_%d", i), fmt.Sprintf("chat_%d", i)
		conn := Conn{name: connName}
		user := NewUser(userName, conn)
		expName := userName
		resName := user.GetName()
		if expName != resName {
			t.Errorf("'GetName()' method FAILED, expected -> %s, got -> %s", expName, resName)
		} else {
			t.Logf("'GetName()' method SUCCEEDED, expected -> %s, got -> %s", expName, resName)
		}
		expConn := conn
		resConn := user.GetConn()
		if expConn != resConn {
			t.Errorf("'GetConn()' method FAILED, expected -> %s, got -> %s", expConn, resConn)
		} else {
			t.Logf("'GetConn()' method SUCCEEDED, expected -> %s, got -> %s", expConn, resConn)
		}
		user.SetRoomName(roomName)
		expRoom := roomName
		resRoom := user.room
		if expRoom != resRoom {
			t.Errorf("'SetRoomName()' method FAILED, expected -> %s, got -> %s", expRoom, resRoom)
		} else {
			t.Logf("'SetRoomName()' method SUCCEEDED, expected -> %s, got -> %s", expRoom, resRoom)
		}
		resRoom, ok := user.GetRoomName()
		if expRoom != resRoom {
			t.Errorf("'GetRoomName()' method FAILED, expected -> %s, got -> %s", expRoom, resRoom)
		} else if !ok {
			t.Errorf("'GetRoomName() method FAILED, expected -> %t, got -> %t", true, false)
		} else {
			t.Logf("'GetRoomName()' method SUCCEEDED, expected -> %s, got -> %s", expRoom, resRoom)
		}
	}

}

//net.Conn interface implementation
type Conn struct {
	name string
}

func (c Conn) Read(b []byte) (int, error) {
	return 0, nil
}

func (c Conn) Write(b []byte) (int, error) {
	return 0, nil
}

func (c Conn) Close() error {
	return nil
}

func (c Conn) LocalAddr() net.Addr {
	return Addr{}
}

func (c Conn) RemoteAddr() net.Addr {
	return Addr{}
}

func (c Conn) SetDeadline(t time.Time) error {
	return nil
}
func (c Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

//net.Addr interface implementation
type Addr struct {
}

func (addr Addr) Network() string {
	return ""
}

func (addr Addr) String() string {
	return ""
}
