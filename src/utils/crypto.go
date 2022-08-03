package utils

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
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

// Fastest way to create a random string given n prefixed by a given prefix (https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go)
var src = rand.NewSource(time.Now().UnixNano())

func RandomString(prefix string, n int) string {
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return fmt.Sprintf("%s%s", prefix, *(*string)(unsafe.Pointer(&b)))
}
