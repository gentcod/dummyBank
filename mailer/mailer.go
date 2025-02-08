package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"sync"
)

// TODO: Implement has a customizable Go package
const (
	gmailSmtpAddress       = "smtp.gmail.com"
	gmailSmtpServerAddress = "smtp.gmail.com:587"
	htmlFilePath           = "../templates/test-mail.html"
	attachmentPath         = "../templates/attachment.pdf"
	subject                = "Account Verification"
	htmlMime               = "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	boundary               = "nextPart"
)

type MailSender interface {
	// sendMail is the main func that holds to the logic for sending emails to a single recipient
	sendMail(recipient Recipient) EmailResult

	// SendEmail is a generic function that helps to send emails to one or more recipients.
	// It handles sending emails to multiple clients concurrently
	SendEmail(recipients ...Recipient) []EmailResult
}

type Mailer struct {
	sender      *GmailSender
	credentials credentials
}

type GmailSender struct {
	smtp.Auth
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

func NewMailer(identity, email, password string) (MailSender, error) {
	sender := newGmailSender(identity, email, password)
	html, err := os.ReadFile(htmlFilePath)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		sender: &sender,
		credentials: credentials{
			identity: identity,
			email:    email,
			password: password,
			template: string(html),
		},
	}, nil
}

func newGmailSender(identity, email, password string) GmailSender {
	auth := smtp.PlainAuth(identity, email, password, gmailSmtpAddress)

	return GmailSender{
		auth,
	}
}

func (mailer *Mailer) SendEmail(recipients ...Recipient) []EmailResult {
	results := make([]EmailResult, 0, len(recipients))
	resultsChan := make(chan EmailResult, len(recipients))

	const maxWorkers = 5
	var wg sync.WaitGroup

	jobs := make(chan Recipient, len(recipients))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for recipient := range jobs {
				result := mailer.sendMail(recipient)
				resultsChan <- result
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

	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (mailer *Mailer) sendMail(recipient Recipient) EmailResult {
	var result EmailResult
	htmlContent, err := generateEmailBody(mailer.credentials.template, recipient.Name, recipient.VerificationLink)
	if err != nil {
		result.Error = err
		return result
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
		mailer.sender.Auth,
		mailer.credentials.email,
		to,
		msg,
	)
	if err != nil {
		result.Error = fmt.Errorf("failed to send message to %s: %v", recipient.Email, err)
		return result
	}

	result.Recipient = recipient
	return result
}
