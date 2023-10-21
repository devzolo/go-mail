package mail

import (
	"testing"
)

func TestMailSend(t *testing.T) {
	es := NewEmailSender(
		WithHost("smtp.hostinger.com"),
		WithPort(465),
		WithUsername(""),
		WithPassword(""),
		WithFrom(""),
	)

	err := es.Send(Email{
		To:      []string{""},
		Subject: "Test1",
		Body:    "Test1",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
