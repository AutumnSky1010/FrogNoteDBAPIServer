package users

// UserId は、ユーザIDを表現する構造体です。
type UserId struct {
	value int
}

// GetValue は、UserIdの値を返却します。
func (id *UserId) GetValue() int {
	return id.value
}

// Equals は、二つのユーザIDの等価性を判定します。等しければtrueを返却します。
func (id *UserId) Equals(otherId *UserId) bool {
	return id.value == otherId.value
}

// NewUserId は、UserId構造体を初期化し、返却します。
func NewUserId(value int) (id *UserId) {
	return &UserId{value: value}
}
