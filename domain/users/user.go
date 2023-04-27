package users

import "errors"

// User は、ユーザ情報を表現する構造体です。
type User struct {
	Id         UserId
	SignInId   SignInId
	ScreenName string
	Password   string
}

// NewUser は、User構造体を初期化し、返却します。screenNameは1文字以上30文字以内、passwordは1文字以上64文字以内です。
func NewUser(id UserId, screenName string, signInId SignInId, password string) (user *User, err error) {
	screenNameLen := len(screenName)
	passwordLen := len(password)
	if screenNameLen > 30 || screenNameLen < 1 {
		return nil, errors.New("'screenName' must be between 1 to 30 characters")
	}
	if passwordLen > 64 || passwordLen < 1 {
		return nil, errors.New("'password' must be between 1 to 64 characters")
	}
	return &User{Id: id, ScreenName: screenName, SignInId: signInId, Password: password}, nil
}
