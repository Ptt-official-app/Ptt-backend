package mail

import (
	"errors"
	"net/smtp"
	"net/url"
	"strings"
)

type smtpProvider struct {
	Mail
	plainAuth smtp.Auth
	smtpURL   *url.URL
}

func createSMTP(url *url.URL) (*smtpProvider, error) {
	username := url.User.Username()
	password, isSet := url.User.Password()

	if !isSet {
		return nil, errors.New("password not set")
	}

	hostSlice := strings.Split(url.Host, ":")
	auth := smtp.PlainAuth("", username, password, hostSlice[0])

	return &smtpProvider{plainAuth: auth, smtpURL: url}, nil
}

func (s *smtpProvider) Send(email, title, userID string, body []byte) error {
	to := []string{email}
	err := smtp.SendMail(s.smtpURL.Host, s.plainAuth, userID+"@ptt-backend.com", to, body)

	if err != nil {
		return err
	}

	return nil
}
