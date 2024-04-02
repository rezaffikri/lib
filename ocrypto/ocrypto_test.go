package ocrypto

import "testing"

const key string = "apqXTEkNmacgXtDC3jmVlo5PvWF8jGlR"

func TestEncryptDecryptImage(t *testing.T) {
	encryptData := EncryptImage([]byte("123456789"), []byte("nkPFOFW8BfV9wVvjLN9apBAlQ0ycUFj4"))
	result := DecryptImage([]byte(encryptData), []byte("nkPFOFW8BfV9wVvjLN9apBAlQ0ycUFj4"))
	if string(result) != "123456789" {
		t.Errorf("EncryptImage() FAILED. Expected %s got %s\n", "123456789", result)
	} else {
		t.Logf("EncryptImage() PASSED. Expected %s got %s\n", "123456789", result)
	}
}

func TestEncryptAESCBC(t *testing.T) {
	plainText := "81$Uhq91?ga"
	encryptData, err := EncryptAESCBC(key, plainText)
	if err != nil {
		t.Errorf("EncryptAESCBC() ERROR. %s\n", err)
	}
	result, err := DecryptAESCBC(key, encryptData)
	if err != nil {
		t.Errorf("EncryptAESCBC() ERROR. %s\n", err)
	}

	if string(result) != plainText {
		t.Errorf("EncryptAESCBC() FAILED. Expected %s got %s\n", plainText, result)
	} else {
		t.Logf("EncryptAESCBC() PASSED. Expected %s got %s\n", plainText, result)
	}
}
