package mail

import (
	"net/smtp"
	"net/url"
	"strings"
)

type smtpProvider struct {
	Mail
	plainAuth smtp.Auth
	url       *url.URL
}

func newSMTPProvider(url *url.URL) (*smtpProvider, error) {
	username := url.User.Username()
	password, _ := url.User.Password()

	hostSlice := strings.Split(url.Host, ":")
	auth := smtp.PlainAuth("", username, password, hostSlice[0])

	return &smtpProvider{plainAuth: auth, url: url}, nil
}

func (s *smtpProvider) Send(from, to, title string, body []byte) error {
	toArr := []string{to}
	err := smtp.SendMail(s.url.Host, s.plainAuth, from, toArr, body)

	if err != nil {
		return err
	}

	return nil
}
