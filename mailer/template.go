package mailer

import (
	"bytes"
	"html/template"
)

type Data struct {
	Name             string
	VerificationLink string
}

// generateEmailBody takes in variables that are used to modify specified `html` template file.
func generateEmailBody(emailTemplate string, data Data) (string, error) {
	tmpl, err := template.New("email").Parse((emailTemplate))
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
