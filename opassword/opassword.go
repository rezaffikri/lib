package opassword

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	number    string = "0123456789"
	special   string = "~=+%^*/()[]{}/!@#$?|"
	lowercase string = "abcdefghijklmnopqrstuvwxyz"
	uppercase string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	space     string = " "
)

var (
	ErrMinLength = errors.New("number of password must be more than than 7 digits")
)

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err == nil {
		rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	} else {
		rand.Seed(time.Now().UnixNano())
	}
}

func HashPassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func GenerateRandom(length int, hasLowerCase, hasUpperCase, hasNumber, hasSpecial, hasSpace bool) (string, error) {
	if length < 6 {
		return "", ErrMinLength
	}
	buf := make([]byte, length)
	var chars string
	if hasLowerCase {
		chars = chars + lowercase
	}
	if hasUpperCase {
		chars = chars + uppercase
	}
	if hasNumber {
		chars = chars + number
	}
	if hasSpecial {
		chars = chars + special
	}
	if hasSpace {
		chars = chars + space
	}

	lowerCaseFound := false
	upperCaseFound := false
	hasNumberFound := false
	hasSpecialFound := false
	hasSpaceFound := false
	for i := 0; i < length; i++ {
		buf[i] = chars[rand.Intn(len(chars))]
		str := string(buf[i])
		if hasLowerCase && !lowerCaseFound && strings.Contains(lowercase, str) {
			lowerCaseFound = true
		}
		if hasUpperCase && !upperCaseFound && strings.Contains(uppercase, str) {
			upperCaseFound = true
		}
		if hasNumber && !hasNumberFound && strings.Contains(number, str) {
			hasNumberFound = true
		}
		if hasSpecial && !hasSpecialFound && strings.Contains(special, str) {
			hasSpecialFound = true
		}
		if hasSpace && !hasSpaceFound && strings.Contains(space, str) {
			hasSpaceFound = true
		}
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	randomIndexLength := length - 1
	randLowerCase := rand.Intn(randomIndexLength)
	if hasLowerCase && !lowerCaseFound {
		buf[randLowerCase] = number[rand.Intn(len(number))]
	}

	randUpperCase := rand.Intn(randomIndexLength)
	if hasUpperCase && !upperCaseFound {
		buf[randUpperCase] = number[rand.Intn(len(number))]
	}

	randNumber := rand.Intn(randomIndexLength)
	if hasNumber && !hasNumberFound {
		buf[randNumber] = number[rand.Intn(len(number))]
	}

	randSpecial := rand.Intn(randomIndexLength)
	if hasSpecial && !hasSpecialFound {
		buf[randSpecial] = number[rand.Intn(len(number))]
	}

	randSpace := rand.Intn(randomIndexLength)
	if hasSpace && !hasSpaceFound {
		buf[randSpace] = number[rand.Intn(len(number))]
	}

	return string(buf), nil
}

type ValidatePasswordConfig struct {
	MinimalPasswordLength   int8
	MaximalPasswordLength   int8
	CheckUpperCasePresent   bool
	CheckLoweCasePresent    bool
	CheckNumberPresent      bool
	CheckSpecialCharPresent bool
}

func NewWeakPasswordConfig() *ValidatePasswordConfig {
	return &ValidatePasswordConfig{
		MinimalPasswordLength:   6,
		MaximalPasswordLength:   32,
		CheckUpperCasePresent:   false,
		CheckLoweCasePresent:    false,
		CheckNumberPresent:      false,
		CheckSpecialCharPresent: false,
	}
}

func NewMediumPasswordConfig() *ValidatePasswordConfig {
	return &ValidatePasswordConfig{
		MinimalPasswordLength:   6,
		MaximalPasswordLength:   32,
		CheckUpperCasePresent:   true,
		CheckLoweCasePresent:    true,
		CheckNumberPresent:      true,
		CheckSpecialCharPresent: false,
	}
}

func NewStrongPasswordConfig() *ValidatePasswordConfig {
	return &ValidatePasswordConfig{
		MinimalPasswordLength:   8,
		MaximalPasswordLength:   32,
		CheckUpperCasePresent:   true,
		CheckLoweCasePresent:    true,
		CheckNumberPresent:      true,
		CheckSpecialCharPresent: true,
	}
}

func ValidatePassword(password string, v *ValidatePasswordConfig) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if v.CheckLoweCasePresent && !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if v.CheckUpperCasePresent && !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if v.CheckNumberPresent && !numberPresent {
		appendError("atleast one numeric character required")
	}
	if v.CheckSpecialCharPresent && !specialCharPresent {
		appendError("special character missing")
	}
	if !(int(v.MinimalPasswordLength) <= passLen && passLen <= int(v.MaximalPasswordLength)) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", v.MinimalPasswordLength, v.MaximalPasswordLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}
