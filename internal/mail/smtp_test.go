package mail

import (
	"reflect"
	"testing"
)

func TestSmtpMail(t *testing.T) {
	provider, err := NewMail("smtp://username:password@mail.smtp.com:587")

	if err != nil {
		t.Errorf("can't get smtp provider")
	}

	providerType := reflect.TypeOf(provider).String()
	if "*mail.smtpProvider" != providerType {
		t.Errorf("provider is not smtp struct type : " + providerType)
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
