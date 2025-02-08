package mailer

import (
	"testing"

	"github.com/gentcod/DummyBank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../test.env")
	require.NoError(t, err)

	email, password := config.MailUser, config.MailPassword
	mailer, err := NewMailer("Dummy Bank", email, password)
	require.NoError(t, err)

	recipient := Recipient{
		Name:             util.RandomOwner(),
		Email:            "oyefuleoluwatayo@gmail.com",
		VerificationLink: "https://github.com/gentcod",
	}

	result := mailer.SendEmail(recipient)
	require.NoError(t, result[0].Error)
}

func TestSendBulkEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../test.env")
	require.NoError(t, err)

	email, password := config.MailUser, config.MailPassword
	mailer, err := NewMailer("Dummy Bank", email, password)
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

	mailResult := mailer.SendEmail(recipients...)
	for _, result := range mailResult {
		require.NoError(t, result.Error)
	}
}
