package oemail

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rezaffikri/lib/oconfig"
)

var config *Config

func init() {
	godotenv.Load(".env.local")
	config = NewConfig(oconfig.GetEnvAsString("SMTP_HOST", ""), oconfig.GetEnvAsString("SMTP_EMAIL", ""), oconfig.GetEnvAsString("SMTP_PASSWORD", ""), oconfig.GetEnvAsString("SMTP_NAME", ""), oconfig.GetEnvAsInt("SMTP_PORT", 0))
}

func TestSendEmailText(t *testing.T) {
	r := NewRequest(config, []string{oconfig.GetEnvAsString("SMTP_TEST_EMAIL", "")}, "TestSendEmail", "test text biasa")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)

	err = r.AttachFile(execPath + "/test.txt")

	if err != nil {
		t.Errorf("SendEmailText() ERROR. %v\n", err)
		return
	}

	ok, err := r.SendEmail()
	if !ok {
		t.Errorf("SendEmailText() ERROR. %v\n", err)
		return
	}

	t.Logf("SendEmailText() SUCCESS. %v\n", ok)
}

func TestSendEmailHTML(t *testing.T) {
	r := NewRequest(config, []string{oconfig.GetEnvAsString("SMTP_TEST_EMAIL", "")}, "TestSendEmail", "")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	path := execPath + "/test_email.min.html"

	templateData := struct {
		Msg string
	}{
		Msg: "Test Parse",
	}

	err = r.ParseTemplate(path, templateData)
	if err != nil {
		t.Errorf("TestSendEmailHTML() ERROR. %v\n", err)
		return
	}

	ok, err := r.SendEmail()
	if !ok {
		t.Errorf("TestSendEmailHTML() ERROR. %v\n", err)
		return
	}

	t.Logf("TestSendEmailHTML() SUCCESS. %v\n", ok)
}

func TestSendEmailHTMLWithAttachment(t *testing.T) {
	r := NewRequest(config, []string{oconfig.GetEnvAsString("SMTP_TEST_EMAIL", "")}, "TestSendEmail", "")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	path := execPath + "/test_email.min.html"

	templateData := struct {
		Msg string
	}{
		Msg: "Test Parse",
	}

	err = r.ParseTemplate(path, templateData)
	if err != nil {
		t.Errorf("TestSendEmailHTMLWithAttachment() ERROR. %v\n", err)
		return
	}

	err = r.AttachFile(execPath + "/test.txt")

	if err != nil {
		t.Errorf("TestSendEmailHTMLWithAttachment() ERROR. %v\n", err)
		return
	}

	ok, err := r.SendEmail()
	if !ok {
		t.Errorf("TestSendEmailHTMLWithAttachment() ERROR. %v\n", err)
		return
	}

	t.Logf("TestSendEmailHTMLWithAttachment() SUCCESS. %v\n", ok)
}

func TestSendEmailHTMLWithMultipleAttachment(t *testing.T) {
	r := NewRequest(config, []string{oconfig.GetEnvAsString("SMTP_TEST_EMAIL", "")}, "TestSendEmail", "")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	path := execPath + "/test_email.min.html"

	templateData := struct {
		Msg string
	}{
		Msg: "Test Parse",
	}

	err = r.ParseTemplate(path, templateData)
	if err != nil {
		t.Errorf("SendEmailHTMLWithMultipleAttachment() ERROR. %v\n", err)
		return
	}

	err = r.AttachFile(execPath + "/test.txt")
	if err != nil {
		t.Errorf("SendEmailHTMLWithMultipleAttachment() ERROR. %v\n", err)
		return
	}

	err = r.AttachFile(execPath + "/test2.txt")
	if err != nil {
		t.Errorf("SendEmailHTMLWithMultipleAttachment() ERROR. %v\n", err)
		return
	}

	ok, err := r.SendEmail()
	if !ok {
		t.Errorf("SendEmailHTMLWithMultipleAttachment() ERROR. %v\n", err)
		return
	}

	t.Logf("SendEmailHTMLWithMultipleAttachment() SUCCESS. %v\n", ok)
}
