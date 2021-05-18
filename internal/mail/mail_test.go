package mail

import (
	"reflect"
	"testing"
)

func TestNewMailProvider_WithSMTP(t *testing.T) {
	provider, err := NewMailProvider("smtp://username@mail.smtp.com:587")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	providerType := reflect.TypeOf(provider).String()
	expectType := reflect.TypeOf(&smtpProvider{}).String()
	if providerType != expectType {
		t.Errorf("expect %s, got %s", expectType, providerType)
	}
}

func TestNewMailProvider_WithMailgun(t *testing.T) {
	provider, err := NewMailProvider("mailgun://host?api_key=xxx&domain=yyy")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	providerType := reflect.TypeOf(provider).String()
	expectType := reflect.TypeOf(&mailgunProvider{}).String()
	if providerType != expectType {
		t.Errorf("expect %s, got %s", expectType, providerType)
	}
}
