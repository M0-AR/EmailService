package main

import (
	"os"
	"testing"
)

func TestSendViaMailgun(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	// Attempt to send an email using the sendViaMailgun method
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")

	// Check if the email sending was not successful
	if !success {
		t.Error("Expected success, but got failure")
	}
}

func TestSendViaMailgunInvalidAPIKey(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: "invalid_api_key",
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid API Key, but got success")
	}
}

func TestSendViaMailgunInvalidDomain(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        "invalid_domain",
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid domain, but got success")
	}
}

func TestSendViaMailgunInvalidEmailAddress(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("invalid_email", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid email address, but got success")
	}
}

func TestSendViaMailgunNoEnvVars(t *testing.T) {
	os.Unsetenv("MAILGUN_PRIVATE_API_KEY")
	os.Unsetenv("MAILGUN_DOMAIN")

	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to missing environment variables, but got success")
	}
}