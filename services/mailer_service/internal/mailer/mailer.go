package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

type MailerService interface {
	SendOTP(email, otp string) error
	SendLoanNotification(email, book string, due time.Time) error
	SendReturnNotification(email, book string) error
}

type mailerService struct {
	email     string
	password  string
	templates map[string]*template.Template
}

// NewMailerService creates a new instance of MailerService
func NewMailerService(email, password string) (MailerService, error) {
	templates := make(map[string]*template.Template)

	// Load all templates
	files := map[string]string{
		"otp":    "./templates/otp_notification.html.gohtml",
		"loan":   "./templates/loan_notification.html.gohtml",
		"return": "./templates/return_notification.html.gohtml",
	}

	for key, path := range files {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
		}
		templates[key] = tmpl
	}

	return &mailerService{
		email:     email,
		password:  password,
		templates: templates,
	}, nil
}

// sendEmail handles the generic email sending logic
func (ms *mailerService) sendEmail(to, subject, templateKey string, data interface{}) error {
	tmpl, exists := ms.templates[templateKey]
	if !exists {
		return fmt.Errorf("template %s not found", templateKey)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	configMessage := gomail.NewMessage()
	configMessage.SetHeader("From", ms.email)
	configMessage.SetHeader("To", to)
	configMessage.SetHeader("Subject", subject)
	configMessage.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("smtp.gmail.com", 587, ms.email, ms.password)

	return dialer.DialAndSend(configMessage)
}

// SendOTP sends an OTP email
func (ms *mailerService) SendOTP(email, otp string) error {
	data := map[string]interface{}{
		"OTP":  otp,
		"Year": time.Now().Year(),
	}
	return ms.sendEmail(email, "Verification Email", "otp", data)
}

// SendLoanNotification sends a loan notification email
func (ms *mailerService) SendLoanNotification(email, book string, due time.Time) error {
	data := map[string]interface{}{
		"Book": book,
		"Due":  due,
	}
	return ms.sendEmail(email, "Loan Notification", "loan", data)
}

// SendReturnNotification sends a return notification email
func (ms *mailerService) SendReturnNotification(email, book string) error {
	data := map[string]interface{}{
		"Book": book,
	}
	return ms.sendEmail(email, "Return Confirmation", "return", data)
}
