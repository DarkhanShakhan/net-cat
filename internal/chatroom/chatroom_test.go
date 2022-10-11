package chatroom

import (
	"fmt"
	"net"
	u "net-cat/internal/userInterface"
	"reflect"
	"testing"
	"time"
)

func TestNewChatroom(t *testing.T) {
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("Room-%d", i)
		resChatroom := NewChatroom(name)
		expChatroom := &Chatroom{name: name, users: map[string]u.User{}}
		if reflect.TypeOf(expChatroom) != reflect.TypeOf(resChatroom) {
			t.Errorf("'NewChatroom' function FAILED, expected %s, got %s", expChatroom, resChatroom)
		} else {
			t.Logf("'NewChatroom' function SUCCEEDED, expected %s, got %s", expChatroom, resChatroom)
		}
		expName := expChatroom.name
		resName := resChatroom.GetChatName()
		if expName != resName {
			t.Errorf("'GetChatName' method FAILED, expected %s, got %s", expName, resName)
		} else {
			t.Logf("'GetChatName' method SUCCEEDED, expected %s, got %s", expName, resName)
		}
		conn := Conn{name: "tcp"}
		username := fmt.Sprintf("user_%d", i)
		user := &User{name: username, conn: conn}
		expChatroom.users[user.name] = user
		expUser := expChatroom.users[username]
		resChatroom.AddUser(user)
		resUser := resChatroom.users[username]
		if expUser != resUser {
			t.Errorf("'AddUser' method FAILED, expected %s, got %s", expUser, resUser)
		} else {
			t.Logf("'AddUser' method SUCCEEDED, expected %s, got %s", expUser, resUser)
		}

	}
}

//to satisfy i.User interface
type User struct {
	name     string
	roomname string
	conn     net.Conn
}

func (u *User) GetName() string {
	return u.name
}
func (u *User) GetConn() net.Conn {
	return u.conn
}

func (u *User) GetRoomName() (string, bool) {
	if u.roomname == "" {
		return "", false
	}
	return u.roomname, true
}

func (u *User) SetRoomName(name string) {
	u.roomname = name
}

//to satisfy net.Conn interface
type Conn struct {
	name string
	log  string
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
