package users_test

import (
	"FrogNote_database/domain/users"
	"testing"
)

func TestNewSignInId(t *testing.T) {
	type testCase struct {
		testName string
		arg      string
	}
	t.Run("有効値", func(t *testing.T) {
		validIdValues := []testCase{
			{testName: "1文字の値", arg: "a"},
			{testName: "30文字の値", arg: "123456789abcdef123456789abcdef"},
		}
		for _, value := range validIdValues {
			t.Run(value.testName, func(t *testing.T) {
				_, err := users.NewSignInId(value.arg)
				if err != nil {
					t.Error(value)
				}
			})
		}
	})

	t.Run("無効値", func(t *testing.T) {
		validIdValues := []testCase{
			{testName: "0文字の値", arg: ""},
			{testName: "31文字の値", arg: "123456789abcdef123456789abcdef1"},
		}
		for _, value := range validIdValues {
			t.Run(value.testName, func(t *testing.T) {
				_, err := users.NewSignInId(value.arg)
				if err == nil {
					t.Error(value)
				}
			})
		}
	})
}

func TestSignInIdEquals(t *testing.T) {
	id1, _ := users.NewSignInId("aiueo")
	id2, _ := users.NewSignInId("aiueo")
	if !id1.Equals(id2) {
		t.Error()
	}
}

func TestSignInIdGetValue(t *testing.T) {
	value := "aiueo"
	id, _ := users.NewSignInId(value)
	if id.GetValue() != value {
		t.Error()
	}
}
