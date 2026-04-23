package mailer

import (
	"crypto/tls"
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
// Port 465 uses implicit TLS (SMTPS); all other ports use STARTTLS via smtp.SendMail.
func (m *Mailer) Send(to, subject, body string) error {
	msg := buildMessage(m.from, to, subject, body)
	addr := m.host + ":" + m.port

	if m.port == "465" {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.host})
		if err != nil {
			return fmt.Errorf("SMTP TLS dial: %w", err)
		}
		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return fmt.Errorf("SMTP client: %w", err)
		}
		defer client.Close()
		if err := client.Auth(smtp.PlainAuth("", m.username, m.password, m.host)); err != nil {
			return fmt.Errorf("SMTP auth: %w", err)
		}
		if err := client.Mail(m.from); err != nil {
			return fmt.Errorf("SMTP MAIL FROM: %w", err)
		}
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("SMTP RCPT TO: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("SMTP DATA: %w", err)
		}
		if _, err := w.Write([]byte(msg)); err != nil {
			return fmt.Errorf("SMTP write: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("SMTP data close: %w", err)
		}
		return client.Quit()
	}

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
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
