//go:build local

package notifications

import "net/smtp"

// smtp.PlainAuth refuses to send your password over an unencrypted connection
// this is a workaround to allow sending emails in a local environment
type unencryptedAuth struct {
	smtp.Auth
}

func NewEmailAuth(identity, username, password, host string) smtp.Auth {
	return &unencryptedAuth{
		smtp.PlainAuth(
			identity,
			username,
			password,
			host,
		),
	}
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
