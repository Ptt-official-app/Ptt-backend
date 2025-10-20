package mail

import (
	"reflect"
	"testing"
)

func TestSmtpMail(t *testing.T) {
	provider, err := NewMailProvider("smtp://username@mail.smtp.com:587")
	if err != nil {
		t.Errorf("can't get smtp provider")
	}

	providerType := reflect.TypeOf(provider).String()
	if providerType != "*mail.smtpProvider" {
		t.Errorf("provider is not smtp struct type : %v", providerType)
	}

	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")

	err = provider.Send("test@example.com", "test", "test", msg)

	if err == nil {
		t.Errorf("send must fail")
	}
}
