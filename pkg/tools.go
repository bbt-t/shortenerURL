package pkg

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"hash/fnv"
	"log"
	"net/url"
	"strings"
)

func HashShortening(s []byte) uint32 {
	/*
		Simple hash func.
		!!! It is NOT a cryptographic hash-func !!!
		return: positive num
	*/
	hash := fnv.New32a()
	if _, err := hash.Write(s); err != nil {
		log.Fatalf("ERROR : %s", err)
	}
	return hash.Sum32()
}

func URLValidation(inpURL string) bool {
	/*
		URL validation.
	*/
	_, err := url.ParseRequestURI(inpURL)
	if err != nil {
		log.Println(err)
	}
	return errors.Is(err, nil)
}

func HostOnly(address string) string {
	/*
		Separating server IP.
		param address: "ip:port"
	*/
	if !strings.Contains(address, ":") {
		return address
	}
	return strings.Split(address, ":")[0]
}

func EncryptPassword(password string) (string, error) {
	/*
		Encrypt the password.
	*/
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func AssertEqualPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return errors.Is(err, nil)
}

func ConvertStrToSlice(strToChange string) []string {
	/*
		[ "a", "b", "c", "d", ...] -> slice [a, b, c, d, ...]
	*/
	temp := strings.ReplaceAll(strToChange, " ", "")
	temp = strings.ReplaceAll(strings.ReplaceAll(temp, "[", ""), "]", "")

	result := strings.Split(strings.ReplaceAll(temp, "\"", ""), ",")

	return result
}
