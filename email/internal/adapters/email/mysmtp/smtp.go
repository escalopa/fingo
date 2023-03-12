package mysmtp

import (
	"fmt"
	"github.com/lordvidex/errs"
	"html/template"
	"net/smtp"
	"time"
)

// Sender is an email sender using smtp
type Sender struct {
	port int
	host string
	from string
	exp  time.Duration
	conn *smtp.Client
	auth smtp.Auth
}

// New creates a new Sender using smtp
func New(opts ...func(*Sender)) (*Sender, error) {
	s := &Sender{}
	// Set options
	for _, opt := range opts {
		opt(s)
	}
	// Create SMTP client
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return nil, errs.B(err).
			Msg("Failed to Dial smtp server, Make sure SMTPSender has port & host set").Err()
	}
	s.conn = conn
	// Create SMTP auth
	s.auth = smtp.PlainAuth("", s.from, "", s.host)
	return s, nil
}

// WithPort sets the port of the smtp server
func WithPort(port int) func(*Sender) {
	return func(s *Sender) {
		s.port = port
	}
}

// WithHost sets the host of the smtp server
func WithHost(host string) func(*Sender) {
	return func(s *Sender) {
		s.host = host
	}
}

// WithFrom sets the from address of the smtp server
func WithFrom(from string) func(*Sender) {
	return func(s *Sender) {
		s.from = from
	}
}

// WithExpiration sets the expiration time of the verification code in the message
func WithExpiration(exp time.Duration) func(*Sender) {
	return func(s *Sender) {
		s.exp = exp
	}
}

// templateData is the data used to create the email message
type templateData struct {
	Code string
	Exp  int
}

// SendVerificationCode sends a verification code to the given email
func (s *Sender) SendVerificationCode(email string, code string) error {
	// Create message
	data := templateData{Code: code, Exp: int(s.exp.Minutes())}
	// TODO: Add a template
	tmp, err := template.New("email").Parse("")
	if err != nil {
		return errs.B(err).Msg("Failed to create template").Err()
	}
	// Execute template
	err = tmp.Execute(nil, data)
	if err != nil {
		return errs.B(err).Msg("Failed to execute template").Err()
	}
	// Send message
	err = smtp.SendMail(s.Url(), nil, s.from, []string{email}, nil)
	if err != nil {
		return errs.B(err).Msg("Failed to send email").Err()
	}
	return nil
}

// Close closes the connection to the smtp server
func (s *Sender) Close() error {
	return s.conn.Close()
}

// Url returns the url of the smtp server, in the form `host:port`
func (s *Sender) Url() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}
