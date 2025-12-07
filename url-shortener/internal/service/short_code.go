package service

import (
	"math/rand"
	"time"
)

const charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

// #nosec
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateShortCode(length int) string {
	answByte := make([]byte, length)
	for i := range answByte {
		answByte[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(answByte)
}
