package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"
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
	decodedHash1, err := hex.DecodeString(hash1)
	if err != nil {
		log.Error("Unable to decode hash: ", err)
		return false
	}
	decodedHash2, err := hex.DecodeString(hash2)
	if err != nil {
		log.Error("Unable to decode hash: ", err)
		return false
	}
	isHashEqual := hmac.Equal(decodedHash1, decodedHash2)

	return isHashEqual
}
