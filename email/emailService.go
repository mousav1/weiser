package email

import (
	"io"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
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

func NewEmailService() *EmailService {
	// Load the configuration file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// Get the SMTP configuration
	smtpConfig := viper.GetStringMap("smtp")
	host := smtpConfig["host"].(string)
	port := smtpConfig["port"].(int)
	username := smtpConfig["username"].(string)
	password := smtpConfig["password"].(string)
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
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
	// Create a new SMTP client
	d := gomail.NewDialer(es.host, es.port, es.username, es.password)
	// Send the message
	if err := d.DialAndSend(m); err != nil {
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
	m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
		return err
	}))
	return nil
}
