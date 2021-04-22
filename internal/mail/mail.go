package mail

import "fmt"

func NewMail(mailDriver string) Mail {
	// TODO: 根據 mailDriver 產生對應 mail
	return &mail{}
}

type mail struct {
	Mail
}

type Mail interface {
	Send(email, title, userID string, body []byte) error
}

func (mail *mail) Send(email, title, userID string, body []byte) error {
	fmt.Printf("call mail send with: %s, %s, %s, %v", email, title, userID, body)
	return nil
}
