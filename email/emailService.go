package email

import (
	"crypto/tls"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer      *gomail.Dialer
	host        string
	port        int
	username    string
	password    string
	attachments []string
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	html        bool
}

func NewEmailService() (*EmailService, error) {
	es := &EmailService{}
	// Load the configuration file
	viper.SetConfigFile("../config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	// Get the SMTP configuration
	smtpConfig := viper.GetStringMap("smtp")
	host, ok := smtpConfig["host"].(string)
	if !ok {
		return nil, errors.New("missing smtp host configuration")
	}
	port, ok := smtpConfig["port"].(int)
	if !ok {
		return nil, errors.New("missing smtp port configuration")
	}
	username, ok := smtpConfig["username"].(string)
	if !ok {
		return nil, errors.New("missing smtp username configuration")
	}
	password, ok := smtpConfig["password"].(string)
	if !ok {
		return nil, errors.New("missing smtp password configuration")
	}
	encryption, ok := smtpConfig["encryption"].(string)
	if !ok {
		encryption = "none"
	}
	mailer, ok := smtpConfig["mailer"].(string)
	if !ok {
		mailer = "smtp"
	}
	es.host = host
	es.port = port
	es.username = username
	es.password = password
	// Create a new SMTP client
	d := gomail.NewDialer(es.host, es.port, es.username, es.password)
	switch encryption {
	case "tls":
		d.TLSConfig = &tls.Config{ServerName: es.host}
	case "ssl":
		d.SSL = true
		d.TLSConfig = &tls.Config{ServerName: es.host}
	}
	switch mailer {
	case "smtp":
		d.LocalName = "localhost"
	}
	es.dialer = d
	return es, nil
}

func (es *EmailService) SetFrom(from string) *EmailService {
	es.from = from
	return es
}

func (es *EmailService) SetTo(to []string) *EmailService {
	es.to = to
	return es
}

func (es *EmailService) SetCc(cc []string) *EmailService {
	es.cc = cc
	return es
}

func (es *EmailService) SetBcc(bcc []string) *EmailService {
	es.bcc = bcc
	return es
}

func (es *EmailService) SetSubject(subject string) *EmailService {
	es.subject = subject
	return es
}

func (es *EmailService) SetBody(body string) *EmailService {
	es.body = body
	return es
}

func (es *EmailService) SetHtml(html bool) *EmailService {
	es.html = html
	return es
}

func (es *EmailService) Attach(file string) *EmailService {
	es.attachments = append(es.attachments, file)
	return es
}

func (es *EmailService) Send() error {
	if es.from == "" {
		return errors.New("missing sender email address")
	}
	if len(es.to) == 0 {
		return errors.New("missing recipient email addresses")
	}
	if es.subject == "" {
		return errors.New("missing email subject")
	}
	if es.body == "" {
		return errors.New("missing email body")
	}

	// Create a new message
	m := gomail.NewMessage()

	// Set the sender and recipient
	m.SetHeader("From", es.from)
	m.SetHeader("To", es.to...)

	// Set the CC and BCC headers
	if len(es.cc) > 0 {
		m.SetHeader("Cc", es.cc...)
	}
	if len(es.bcc) > 0 {
		m.SetHeader("Bcc", es.bcc...)
	}

	// Set the subject and body
	m.SetHeader("Subject", es.subject)
	if es.html {
		m.SetBody("text/html", es.body)
	} else {
		m.SetBody("text/plain", es.body)
	}

	// Add attachments
	for _, file := range es.attachments {
		if err := es.addAttachment(m, file); err != nil {
			return err
		}
	}

	// Send the message using the existing SMTP client
	if err := es.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (es *EmailService) addAttachment(m *gomail.Message, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add the attachment to the message
	_, fileName := filepath.Split(filename)
	m.Attach(fileName,
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := file.Seek(0, 0)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, file)
			return err
		}),
	)

	return nil
}
