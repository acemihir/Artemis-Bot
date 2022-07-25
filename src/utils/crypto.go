package utils

import (
	"crypto/aes"
	"encoding/hex"
)

func EncryptAES(key []byte, plaintext string) string {
	c, ex := aes.NewCipher(key)
	if ex != nil {
		return ""
	}

	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) string {
	ciptext, _ := hex.DecodeString(ct)

	c, ex := aes.NewCipher(key)
	if ex != nil {
		return ""
	}

	pt := make([]byte, len(ciptext))
	c.Decrypt(pt, ciptext)

	s := string(pt[:])
	return s
}
