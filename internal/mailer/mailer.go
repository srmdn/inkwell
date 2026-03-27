package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Mailer sends email via SMTP.
type Mailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func New(host, port, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Send sends a plain-text email to a single recipient.
func (m *Mailer) Send(to, subject, body string) error {
	addr := m.host + ":" + m.port
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	msg := buildMessage(m.from, to, subject, body)
	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
}

// SendMany sends the same email to multiple recipients one by one.
// Returns a combined error listing any failures.
func (m *Mailer) SendMany(recipients []string, subject, body string) error {
	var errs []string
	for _, to := range recipients {
		if err := m.Send(to, subject, body); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", to, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("send failures: %s", strings.Join(errs, "; "))
	}
	return nil
}

func buildMessage(from, to, subject, body string) string {
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	return sb.String()
}
