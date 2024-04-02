package opassword

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	plainText := "123"
	hashText := HashPassword([]byte(plainText))
	result := ComparePasswords(hashText, []byte(plainText))
	if result == false {
		t.Errorf("HashPassword() FAILED. Expected %v got %v\n", true, result)
	} else {
		t.Logf("HashPassword() PASSED. Expected %v got %v\n", true, result)
	}
}

func TestGeneratePassword(t *testing.T) {
	generatedPassword, err := GenerateRandom(8, true, true, true, true, true)
	if err != nil {
		t.Errorf("GeneratePassword() ERROR. %v\n", err)
	}
	passCfg := NewStrongPasswordConfig()
	err = ValidatePassword(generatedPassword, passCfg)
	if err != nil {
		t.Errorf("GeneratePassword() ERROR. %v\n", err)
	} else {
		t.Logf("GeneratePassword() PASSED")
	}
}
