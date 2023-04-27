package security_test

import (
	"FrogNote_database/domain/users"
	"FrogNote_database/infrastructure/security"
	"testing"
)

func TestGetUserId(t *testing.T) {
	tokens := security.Tokens{}
	userId := users.NewUserId(1)
	token := tokens.GenereteToken(userId)
	userIdFromTokens, ok := tokens.GetUserId(token)

	if !ok || !userId.Equals(userIdFromTokens) {
		t.Error()
	}
}

func TestInvalidate(t *testing.T) {
	tokens := security.Tokens{}
	userId := users.NewUserId(1)
	token := tokens.GenereteToken(userId)

	// 無効化
	tokens.Invalidate(token)

	_, ok := tokens.GetUserId(token)
	// 無効化したので取得できたらアウト
	if ok {
		t.Error()
	}
}
