package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hasher struct {
	key string
}

func New(key string) *Hasher {
	return &Hasher{
		key: key,
	}
}

func (h *Hasher) Hash(str string) string {
	if h.key == "" {
		return ""
	}

	hmacHash := hmac.New(sha256.New, []byte(h.key))
	hmacHash.Write([]byte(str))
	hash := hmacHash.Sum(nil)

	return fmt.Sprintf("%x", hash)
}

func (h *Hasher) IsEqual(hash1 string, hash2 string) bool {
	decodedHash1, _ := hex.DecodeString(hash1)
	decodedHash2, _ := hex.DecodeString(hash2)
	isHashEqual := hmac.Equal(decodedHash1, decodedHash2)

	return isHashEqual
}
