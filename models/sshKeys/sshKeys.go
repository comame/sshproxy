package sshKeys

import (
	"fmt"
	"strings"
)

type Option struct {
	PermitOpen        []address
	NoAgentForwarding bool
	NoUserRc          bool
	NoX11Forwarding   bool
	NoPty             bool
	Command           []string
}

func (o Option) String() string {
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

// 渡された SSH 公開鍵とオプションから、authorized_keys の行を生成する。
// `sshKey` は有効な公開鍵のフォーマットであることを想定する
func CreateAuthorizedKeyLine(sshKey string, opt Option) string {
	if !IsValidSSHPublicKey(sshKey) {
		panic("got malformed ssh public key")
	}

	optStr := opt.String()
	if optStr == "" {
		return sshKey
	}
	return optStr + " " + sshKey
}

var (
	keyTypeRSA     = "ssh-rsa"
	keyTypeED25519 = "ssh-ed25519"
)

func isKeyType(str string) bool {
	return str == keyTypeRSA || str == keyTypeED25519
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

// 引数が有効な SSH 公開鍵であることを検証する
func IsValidSSHPublicKey(str string) bool {
	s := strings.Split(str, " ")
	if len(s) != 3 {
		return false
	}

	if !isKeyType(s[0]) {
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
