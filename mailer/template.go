package mailer

import (
	"bytes"
	"html/template"
)

// generateEmailBody takes in variables that are used to modify specified `html` template file.
func generateEmailBody(emailTemplate, name, verificationLink string) (string, error) {
	tmpl, err := template.New("email").Parse((emailTemplate))
	if err != nil {
		return "", err
	}

	data := struct {
		Name             string
		VerificationLink string
	}{
		Name:             name,
		VerificationLink: verificationLink,
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
