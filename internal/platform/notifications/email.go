package notifications

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SMTPService struct {
	Host      string
	Port      string
	Username  string
	Password  string
	logger    *zap.Logger
	templates map[string]string
}

//go:embed template/summary.html
var summaryTemplate string

func NewSMTPService(config *SMTPConfig, logger *zap.Logger) *SMTPService {
	return &SMTPService{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		logger:   logger,
		templates: map[string]string{
			"summary": summaryTemplate,
		},
	}
}

func (s *SMTPService) Send(to, subject, templateName string, variables map[string]interface{}) error {
	emailTemplate, ok := s.templates[templateName]
	if !ok {
		s.logger.Error("template not found", zap.String("template", templateName))
		return fmt.Errorf("template %s not found", templateName)
	}
	s.logger.Info("Sending email", zap.String("to", to), zap.String("subject", subject))
	renderedContent, err := s.parseTemplate(emailTemplate, variables)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Build the email message
	message := s.buildMessage(s.Username, to, subject, renderedContent)

	// SMTP address and authentication
	smtpAddr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	s.logger.Info("Sending email", zap.String("SmtpAddr", smtpAddr))
	auth := NewEmailAuth("", s.Username, s.Password, s.Host)

	// Send the email
	err = smtp.SendMail(smtpAddr, auth, s.Username, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	fmt.Printf("Email sent to %s with subject: %s\n", to, subject)
	return nil
}

// parseTemplate replaces variables in the template string
func (s *SMTPService) parseTemplate(templateString string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := tmpl.Execute(&builder, variables); err != nil {
		return "", err
	}

	return builder.String(), nil
}

// buildMessage constructs the email message body
func (s *SMTPService) buildMessage(from, to, subject, body string) string {
	return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body)
}
