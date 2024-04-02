package oemail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SMTP_HOST     string
	SMTP_PORT     int
	SMTP_EMAIL    string
	SMTP_PASSWORD string
	SMTP_NAME     string
}

// Request struct
type EmailRequest struct {
	config      *Config
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	isBodyHTML  bool
	attachments map[string][]byte
}

func NewConfig(SMTP_HOST, SMTP_EMAIL, SMTP_PASSWORD, SMTP_NAME string, SMTP_PORT int) *Config {
	return &Config{
		SMTP_HOST:     SMTP_HOST,
		SMTP_PORT:     SMTP_PORT,
		SMTP_EMAIL:    SMTP_EMAIL,
		SMTP_PASSWORD: SMTP_PASSWORD,
		SMTP_NAME:     SMTP_NAME,
	}
}

func NewRequest(config *Config, to []string, subject, body string) *EmailRequest {
	return &EmailRequest{
		config:      config,
		to:          to,
		subject:     subject,
		body:        body,
		attachments: make(map[string][]byte),
	}
}

func (r *EmailRequest) SendEmail() (bool, error) {
	addr := fmt.Sprintf("%s:%d", r.config.SMTP_HOST, r.config.SMTP_PORT)

	auth := smtp.PlainAuth("", r.config.SMTP_EMAIL, r.config.SMTP_PASSWORD, r.config.SMTP_HOST)

	if err := smtp.SendMail(addr, auth, r.config.SMTP_EMAIL, r.to, r.ToBytes()); err != nil {
		return false, err
	}
	return true, nil
}

func (r *EmailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	r.isBodyHTML = true
	return nil
}

func (m *EmailRequest) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.attachments[fileName] = b
	return nil
}

func (m *EmailRequest) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.to, ",")))
	if len(m.cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.cc, ",")))
	}

	if len(m.bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.bcc, ",")))
	}

	buf.WriteString("MIME-Version: 1.0;\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	if m.isBodyHTML {
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n\n")
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.body)
	if withAttachments {
		for k, v := range m.attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
