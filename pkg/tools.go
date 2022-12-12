package pkg

import (
	"hash/fnv"
	"log"
	"net/url"
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
	return nil == err
}
