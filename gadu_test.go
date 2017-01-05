package gadu

import "testing"

func TestVersion(t *testing.T) {
	version := Version()
	if version != "1.12.1" {
		t.Errorf("Invalid version: %s", version)
	}
}

func TestGGSession(t *testing.T) {
	session := NewGGSession()
	defer session.Close()
	_ = session
}

func TestGGSessionLogin(t *testing.T) {
	session := NewGGSession()
	defer session.Close()
	session.Uin = 1234
	session.Password = "password"

	if e := session.Login(); e != AccessDeniedError {
		t.Fatalf("Unable to login: %s", e.Error())
	}
}

func TestGGSessionSendMessage(t *testing.T) {
	session := NewGGSession()
	defer session.Close()
	session.Uin = 1234
	session.Password = "password"
	if e := session.Login(); e != AccessDeniedError {
		t.Fatalf("Unable to login: %s", e.Error())
	}
	session.SendMessage(123456789, "Hello, world!")
}
