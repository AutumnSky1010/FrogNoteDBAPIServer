package users

import "errors"

// SignInId は、サインインIDを表現する構造体です。
type SignInId struct {
	value string
}

// GetValue は、サインインIDの値を返却します。
func (id *SignInId) GetValue() string {
	return id.value
}

// Equals は、二つのサインインIDの等価性を比較します。等しければtrueを返却します。
func (id *SignInId) Equals(otherId *SignInId) bool {
	return id.value == otherId.value
}

// NewSignInId は、SignInIdを初期化し、返却します。引数のvalueは、1文字以上30文字以内です。
func NewSignInId(value string) (id *SignInId, err error) {
	valueLen := len(value)
	if valueLen > 30 || valueLen < 1 {
		return nil, errors.New("'value' must be between 1 to 30 characters")
	}
	return &SignInId{value: value}, nil
}
