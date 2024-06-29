package sshKeys

import "testing"

func TestIsAllowedChar(t *testing.T) {
	if isAllowedChar('a', "bcd") {
		t.FailNow()
	}
	if !isAllowedChar('a', "abcd") {
		t.FailNow()
	}
}

func TestIsAllowedString(t *testing.T) {
	if !isAllowedString("HelloWorld314", asciiChars) {
		t.FailNow()
	}
	if isAllowedString("Hello, World314", asciiChars) {
		t.FailNow()
	}
}

func TestIsSSHKey(t *testing.T) {
	if !IsValidSSHPublicKey("ssh-rsa abcde test@example.com") {
		t.FailNow()
	}

	if IsValidSSHPublicKey("agent-forwarding ssh-rsa abcde test@example.com") {
		t.FailNow()
	}
	if IsValidSSHPublicKey("agent-forwarding,command=\"/bin/bash\" ssh-rsa abcde test@example.com") {
		t.FailNow()
	}
}

func TestCreateAuthorizedKeyLine_emptyOption(t *testing.T) {
	opt := Option{
		PermitOpen:        nil,
		NoAgentForwarding: false,
		NoUserRc:          false,
		NoX11Forwarding:   false,
		NoPty:             false,
		Command:           nil,
	}
	k := "ssh-rsa abcde== user@example.com"

	expect := "ssh-rsa abcde== user@example.com"

	if got := CreateAuthorizedKeyLine(k, opt); got != expect {
		t.Errorf("フォーマットが違う: %s", got)
	}
}

func TestCreateAuthorizedKeyLine_withOption(t *testing.T) {
	opt := Option{
		PermitOpen: []address{
			{Hostname: "host1.example.com", Port: 8080},
			{Hostname: "host2.example.com", Port: 8081},
		},
		NoAgentForwarding: true,
		NoUserRc:          true,
		NoX11Forwarding:   true,
		NoPty:             true,
		Command:           nil,
	}
	k := "ssh-rsa abcde== user@example.com"

	expect := "permitopen=\"host1.example.com:8080\",permitopen=\"host2.example.com:8081\",no-agent-forwarding,no-user-rc,no-x11-forwarding,no-pty ssh-rsa abcde== user@example.com"

	if got := CreateAuthorizedKeyLine(k, opt); got != expect {
		t.Errorf("フォーマットが違う: %s", got)
	}
}
