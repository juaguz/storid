//go:build !local

package notifications

import "net/smtp"

func NewEmailAuth(identity, username, password, host string) smtp.Auth {
	return smtp.PlainAuth(
		identity,
		username,
		password,
		host)

}
