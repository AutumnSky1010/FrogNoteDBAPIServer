package security_test

import (
	"FrogNote_database/infrastructure/security"
	"testing"
)

func TestGetHash(t *testing.T) {
	str := "a"
	hash := security.Hash{}
	hashStr := hash.GetHash(str)

	// 64桁ならば、文字列が加工されているはず。
	t.Run("文字数が64桁になっているかのテスト", func(t *testing.T) {
		if len(hashStr) != 64 {
			t.Error()
		}
	})

	t.Run("同じ文字列をハッシュ化した際、同じ結果になっているかのテスト", func(t *testing.T) {
		hashStr2 := hash.GetHash(str)
		if hashStr != hashStr2 {
			t.Error()
		}
	})
}
