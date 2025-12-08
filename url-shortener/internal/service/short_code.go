package service

import "crypto/sha256"

const charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func GenerateShortCode(longURL, salt string, length int) string {
	if length <= 0 {
		return ""
	}

	data := []byte(longURL + salt)

	hash := sha256.Sum256(data)

	answByte := make([]byte, length)

	for i := range length {
		ambit := (i * 4) % (len(hash) - 3)

		var sum int
		for _, val := range hash[ambit : ambit+4] {
			sum += int(val)
		}

		answByte[i] = charset[sum%len(charset)]
	}

	return string(answByte)
}
