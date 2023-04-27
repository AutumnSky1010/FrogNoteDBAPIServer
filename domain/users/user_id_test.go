package users_test

import (
	"FrogNote_database/domain/users"
	"testing"
)

func TestUserIdGetValue(t *testing.T) {
	value := 1
	id := users.NewUserId(value)
	if id.GetValue() != value {
		t.Error()
	}
}

func TestUserIdEquals(t *testing.T) {
	value := 1
	id1 := users.NewUserId(value)
	id2 := users.NewUserId(value)

	if !id1.Equals(id2) {
		t.Error()
	}
}
