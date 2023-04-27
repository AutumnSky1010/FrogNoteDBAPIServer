package security

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash は、文字列をハッシュ化する機能を提供する構造体です。
type Hash struct{}

// GetHash は、与えられた文字列をSHA-256で、10,000回ハッシュ化します。
func (h *Hash) GetHash(str string) string {
	sha := sha256.New()
	result := str
	for i := 0; i < 10000; i++ {
		sha.Write([]byte(result))
		strBytes := sha.Sum(nil)
		result = hex.EncodeToString(strBytes)
	}
	return result
}
