package service

import (
	"crypto/sha1"
	"fmt"
)

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func compareHashAndPassword(hash, password, salt string) bool {
	return generatePasswordHash(password, salt) == hash
}
