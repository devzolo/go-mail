package mail

import (
	"crypto/tls"
	"errors"

	"gopkg.in/gomail.v2"
)

const (
	ErrNoRecipients = "no recipients specified"
)

// EmailConfig contains configuration settings for the email sender.
type EmailConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	TLSConfig *tls.Config
	SSLMode   bool
}

// EmailSender is responsible for sending emails.
type EmailSender struct {
	config EmailConfig
}

// Email contains details of the email to be sent.
type Email struct {
	To          []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string
}

func (email *Email) SetHTMLBody(body string) {
	email.Body = body
	email.IsHTML = true
}

func (email *Email) SetTextBody(body string) {
	email.Body = body
	email.IsHTML = false
}

// EmailSenderOptionSetter is a function type that sets options for EmailSender.
type EmailSenderOptionSetter func(*EmailSender)

// NewEmailSender initializes a new EmailSender with the provided options.
func NewEmailSender(options ...EmailSenderOptionSetter) *EmailSender {
	es := &EmailSender{}

	for _, option := range options {
		option(es)
	}

	return es
}

// isValid checks if the email is valid before sending.
func (email Email) isValid() error {
	if len(email.To) == 0 {
		return errors.New(ErrNoRecipients)
	}
	return nil
}

// Send sends the provided email.
func (es *EmailSender) Send(email Email) error {
	if err := email.isValid(); err != nil {
		return err
	}

	m := es.createMessage(email)

	return es.dialAndSend(m)
}

// createMessage creates the email message from the provided email details.
func (es *EmailSender) createMessage(email Email) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", es.config.From)
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)

	contentType := "text/plain"
	if email.IsHTML {
		contentType = "text/html"
	}
	m.SetBody(contentType, email.Body)

	for _, attachment := range email.Attachments {
		m.Attach(attachment)
	}

	return m
}

func (es *EmailSender) dialAndSend(m *gomail.Message) error {
	d := gomail.NewDialer(es.config.Host, es.config.Port, es.config.Username, es.config.Password)
	d.TLSConfig = es.config.TLSConfig

	if es.config.SSLMode {
		d.SSL = true
	}

	return d.DialAndSend(m)
}

// WithHost sets the host for the EmailSender.
func WithHost(host string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Host = host
	}
}

// WithPort sets the port for the EmailSender.
func WithPort(port int) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Port = port
	}
}

// WithUsername sets the username for the EmailSender.
func WithUsername(username string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Username = username
	}
}

// WithPassword sets the password for the EmailSender.
func WithPassword(password string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.Password = password
	}
}

// WithFrom sets the sender address for the EmailSender.
func WithFrom(from string) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.From = from
	}
}

// WithTLSConfig sets the TLS configuration for the EmailSender.
func WithTLSConfig(config *tls.Config) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.TLSConfig = config
	}
}

// WithSSLMode sets the SSL mode for the EmailSender.
func WithSSLMode(sslMode bool) EmailSenderOptionSetter {
	return func(es *EmailSender) {
		es.config.SSLMode = sslMode
	}
}
