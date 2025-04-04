package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"sync"
)

// TODO: Implement as a customizable Go package
const (
	gmailSmtpAddress       = "smtp.gmail.com"
	gmailSmtpServerAddress = "smtp.gmail.com:587"
	attachmentPath         = "../templates/attachment.pdf"
	subject                = "Account Verification"
	htmlMime               = "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	boundary               = "nextPart"
)

type MailSender interface {
	// sendMail is the main func that holds to the logic for sending emails to a single recipient
	sendMail(recipient Recipient) error

	// SendEmail is a generic function that helps to send emails to one or more recipients.
	// It handles sending emails to multiple clients concurrently
	SendEmail(recipients ...Recipient) error
}

type GmailSender struct {
	auth        smtp.Auth
	credentials credentials
}

type credentials struct {
	identity, email, password, template string
}

type EmailResult struct {
	Recipient Recipient
	Error     error
}

type Recipient struct {
	Name             string
	Email            string
	VerificationLink string
}

func NewMailer(identity, email, password string, html []byte) (MailSender, error) {
	auth := smtp.PlainAuth(identity, email, password, gmailSmtpAddress)

	return &GmailSender{
		auth: auth,
		credentials: credentials{
			identity: identity,
			email:    email,
			password: password,
			template: string(html),
		},
	}, nil
}

func (mailer *GmailSender) SendEmail(recipients ...Recipient) error {
	resultsChan := make(chan error, len(recipients))

	const maxWorkers = 5
	var wg sync.WaitGroup

	jobs := make(chan Recipient, len(recipients))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for recipient := range jobs {
				err := mailer.sendMail(recipient)
				resultsChan <- err
			}
		}()
	}

	for _, recipient := range recipients {
		jobs <- recipient
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for err := range resultsChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (mailer *GmailSender) sendMail(recipient Recipient) error {
	htmlContent, err := generateEmailBody(
		mailer.credentials.template,
		Data{
			Name:             recipient.Name,
			VerificationLink: recipient.VerificationLink,
		},
	)
	if err != nil {
		return err
	}

	headers := fmt.Sprintf("From: %s\r\n", mailer.credentials.identity)
	headers += fmt.Sprintf("To: %s\r\n", recipient.Email)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n"

	var body bytes.Buffer
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	body.WriteString("\r\n" + string(htmlContent) + "\r\n\r\n")

	to := []string{recipient.Email}
	msg := []byte(headers + body.String())

	err = smtp.SendMail(
		gmailSmtpServerAddress,
		mailer.auth,
		mailer.credentials.email,
		to,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send message to %s: %v", recipient.Email, err)
	}

	return nil
}
