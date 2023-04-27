package security

import (
	"FrogNote_database/domain/users"

	"github.com/google/uuid"
)

var (
	// tokenUserPairs は、トークンとユーザIDを紐づけるためのマップです。
	tokenUserPairs map[string]*users.UserId = make(map[string]*users.UserId)
)

// Tokens は、トークンの生成、無効化を行う構造体です。
type Tokens struct{}

// GenerateToken は、新しいトークンを生成します。
func (tokens *Tokens) GenereteToken(userId *users.UserId) (token string) {
	notGenereted := true
	for notGenereted {
		uuid, err := uuid.NewRandom()
		if err != nil {
			continue
		}

		token = uuid.String()
		// すでに同じトークンが存在していた場合、やりなおす。（uuidなので基本的に重複しないが、セキュリティ機能なのでしっかり重複チェックを行っている。）
		if _, exists := tokens.GetUserId(token); exists {
			continue
		}
		// 存在していない
		notGenereted = false
	}
	tokens.set(token, userId)
	return
}

// Invalidate は、トークンを無効化します。
func (tokens *Tokens) Invalidate(token string) {
	delete(tokenUserPairs, token)
}

// GetUserId は、トークンに紐づけられたユーザを取得します。
func (tokens *Tokens) GetUserId(token string) (userId *users.UserId, ok bool) {
	userId, ok = tokenUserPairs[token]
	return
}

// set は、トークンとユーザIDを紐づけます。
func (tokens *Tokens) set(token string, userId *users.UserId) {
	tokenUserPairs[token] = userId
}
