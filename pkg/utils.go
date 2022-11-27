package pkg

import (
	"hash/fnv"
	"log"
)

func HashShortening(s []byte) uint32 {
	hash := fnv.New32a()
	if _, err := hash.Write(s); err != nil {
		log.Fatalf("ERROR : %s", err)
	}
	return hash.Sum32()
}
