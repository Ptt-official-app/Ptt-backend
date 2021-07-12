package mail

import (
	"net/smtp"
	"net/url"
)

type smtpProvider struct {
	Mail
	plainAuth smtp.Auth
	url       *url.URL
}

func newSMTPProvider(url *url.URL) *smtpProvider {
	username := url.User.Username()
	password, _ := url.User.Password()

	hostName := url.Hostname()
	auth := smtp.PlainAuth("", username, password, hostName)

	return &smtpProvider{plainAuth: auth, url: url}
}

func (s *smtpProvider) Send(from, to, title string, body []byte) error {
	toArr := []string{to}
	err := smtp.SendMail(s.url.Host, s.plainAuth, from, toArr, body)
	if err != nil {
		return err
	}

	return nil
}
