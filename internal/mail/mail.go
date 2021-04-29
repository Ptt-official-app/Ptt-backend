package mail

import (
	"fmt"
	"net/url"
)

//stmp example - smtp://username:password@mail.smtp.com:25
func NewMailProvider(mailDriver string) (Mail, error) {
	urlStruct, err := url.Parse(mailDriver)
	if err != nil {
		return nil, err
	}

	switch urlStruct.Scheme {
	case "smtp":
		provider, err := createSMTP(urlStruct)

		if err != nil {
			return nil, err
		}

		return provider, nil
	}

	return &mail{}, nil
}

type mail struct {
	Mail
}

type Mail interface {
	Send(email, title, userID string, body []byte) error
}

func (mail *mail) Send(from, to, title string, body []byte) error {
	fmt.Printf("call mail send with: %s, %s, %v", from, title, body)
	return nil
}
