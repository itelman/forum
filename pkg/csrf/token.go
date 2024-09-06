package csrf

import (
	"crypto/rand"
	"encoding/base64"
)

func NewToken() (string, error) {
	tokenBytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}
