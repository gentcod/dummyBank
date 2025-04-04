package mailer

import (
	"os"
	"testing"

	"github.com/gentcod/DummyBank/util"
	"github.com/stretchr/testify/require"
)

const (
	htmlFilePath = "../templates/test-mail.html"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../test.env")
	require.NoError(t, err)

	html, err := os.ReadFile(htmlFilePath)
	require.NoError(t, err)

	mailer, err := NewMailer("Dummy Bank", config.MailUser, config.MailPassword, html)
	require.NoError(t, err)

	recipient := Recipient{
		Name:             util.RandomOwner(),
		Email:            "oyefuleoluwatayo@gmail.com",
		VerificationLink: "https://github.com/gentcod",
	}

	err = mailer.SendEmail(recipient)
	require.NoError(t, err)
}

func TestSendBulkEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../test.env")
	require.NoError(t, err)

	html, err := os.ReadFile(htmlFilePath)
	require.NoError(t, err)

	mailer, err := NewMailer("Dummy Bank", config.MailUser, config.MailPassword, html)
	require.NoError(t, err)

	recipients := []Recipient{
		{
			Name:             util.RandomOwner(),
			Email:            "oyefuleoluwatayo@gmail.com",
			VerificationLink: "https://github.com/gentcod",
		},
		{
			Name:             util.RandomOwner(),
			Email:            "drelanorgent@gmail.com",
			VerificationLink: "https://github.com/gentcod",
		},
		{
			Name:             util.RandomOwner(),
			Email:            "oye.grox@gmail.com",
			VerificationLink: "https://github.com/gentcod",
		},
	}

	err = mailer.SendEmail(recipients...)
	require.NoError(t, err)
}
