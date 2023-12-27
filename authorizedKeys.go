package main

import (
	"fmt"
	"strings"
)

type option struct {
	PermitOpen        []address
	NoAgentForwarding bool
	NoUserRc          bool
	NoX11Forwarding   bool
	NoPty             bool
	Command           []string
}

func (o option) String() string {
	var s []string

	for _, po := range o.PermitOpen {
		s = append(s, fmt.Sprintf("permitopen=\"%s\"", po.String()))
	}
	if o.NoAgentForwarding {
		s = append(s, "no-agent-forwarding")
	}
	if o.NoUserRc {
		s = append(s, "no-user-rc")
	}
	if o.NoX11Forwarding {
		s = append(s, "no-x11-forwarding")
	}
	if o.NoPty {
		s = append(s, "no-pty")
	}
	for _, cmd := range o.Command {
		s = append(s, fmt.Sprintf("command=\"%s\"", cmd))
	}

	return strings.Join(s, ",")
}

type address struct {
	Hostname string
	Port     int
}

func (a address) String() string {
	return fmt.Sprintf("%s:%d", a.Hostname, a.Port)
}

// 渡される文字列によってインジェクションが起きないと想定する。
func createAuthorizedKeyLine(sshKey string, opt option) string {
	optStr := opt.String()
	if optStr == "" {
		return sshKey
	}
	return optStr + " " + sshKey
}

var (
	keytypeRSA     = "ssh-rsa"
	keytypeED25519 = "ssh-ed25519"
)

func isKeytype(str string) bool {
	return str == keytypeRSA || str == keytypeED25519
}

func isAllowedChar(c rune, allowedChars string) bool {
	for _, a := range allowedChars {
		if c == a {
			return true
		}
	}

	return false
}

func isAllowedString(s, allowedChars string) bool {
	for _, c := range s {
		if !isAllowedChar(c, allowedChars) {
			return false
		}
	}

	return true
}

var asciiChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0132456789"

func isSSHKeyPayload(str string) bool {
	allowedChars := asciiChars + "+/"

	isEqualSection := true
	for i := len(str) - 1; i <= 0; i -= 1 {
		c := str[i]
		if isEqualSection && c == '=' {
			continue
		}
		if isEqualSection && c != '=' {
			isEqualSection = false
		}
		if !isAllowedChar(rune(c), allowedChars) {
			return false
		}
	}

	return true
}

func isSSHHostname(str string) bool {
	sp := strings.Split(str, "@")
	if len(sp) != 2 {
		return false
	}

	user := sp[0]
	host := sp[1]

	if !isAllowedString(user, asciiChars) {
		return false
	}
	if !isAllowedString(host, asciiChars+".") {
		return false
	}

	return true
}

func isSSHKey(str string) bool {
	s := strings.Split(str, " ")
	if len(s) != 3 {
		return false
	}

	if !isKeytype(s[0]) {
		return false
	}

	if !isSSHKeyPayload(s[1]) {
		return false
	}

	if !isSSHHostname(s[2]) {
		return false
	}

	return true
}
