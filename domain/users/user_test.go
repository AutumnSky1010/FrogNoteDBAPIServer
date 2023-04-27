package users_test

import (
	"FrogNote_database/domain/users"
	"testing"
)

func TestNewUser(t *testing.T) {
	type testCase struct {
		testName string
		arg      string
	}
	signInId, _ := users.NewSignInId("aiueo")
	userId := *users.NewUserId(1)
	t.Run("正常", func(t *testing.T) {
		validPasswords := []testCase{
			{testName: "1文字のパスワード", arg: "a"},
			{testName: "64文字のパスワード", arg: "1234567812345678123456781234567812345678123456781234567812345678"},
		}
		validScreenNames := []testCase{
			{testName: "1文字のスクリーンネーム", arg: "a"},
			{testName: "30文字のスクリーンネーム", arg: "123456789abcdef123456789abcdef"},
		}
		for i := 0; i < len(validPasswords); i++ {
			for j := 0; j < len(validScreenNames); j++ {
				t.Run(validPasswords[i].testName+validScreenNames[j].testName, func(t *testing.T) {
					_, err := users.NewUser(userId, validScreenNames[j].arg, *signInId, validPasswords[i].arg)
					if err != nil {
						t.Error()
						return
					}
				})
			}
		}
	})
	t.Run("異常", func(t *testing.T) {
		invalidPasswords := []testCase{
			{testName: "0文字のパスワード", arg: ""},
			{testName: "65文字のパスワード", arg: "12345678123456781234567812345678123456781234567812345678123456781"},
		}
		invalidScreenNames := []testCase{
			{testName: "0文字のスクリーンネーム", arg: ""},
			{testName: "31文字のスクリーンネーム", arg: "123456789abcdef123456789abcdef1"},
		}
		for i := 0; i < len(invalidPasswords); i++ {
			for j := 0; j < len(invalidScreenNames); j++ {
				t.Run(invalidPasswords[i].testName+invalidScreenNames[j].testName, func(t *testing.T) {
					_, err := users.NewUser(userId, invalidScreenNames[j].arg, *signInId, invalidPasswords[i].arg)
					if err == nil {
						t.Error()
						return
					}
				})
			}
		}
	})
}
