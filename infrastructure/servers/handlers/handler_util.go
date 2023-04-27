package handlers

import (
	domainUsers "FrogNote_database/domain/users"
	"FrogNote_database/infrastructure/security"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// IsNotJsonReq は、リクエストボディがJSONであることと、想定されたHTTPメソッドかを判定します。異なった場合はtrueが返却されます。
func IsNotJsonReq(req *http.Request, httpMethod string) bool {
	return req.Method != httpMethod || req.Header.Get("Content-Type") != "application/json"
}

// GetUserId は、リクエストのAuthorizationヘッダをもとにユーザIDを取得します。
func GetUserId(req *http.Request) (id *domainUsers.UserId, ok bool) {
	tokens := security.Tokens{}
	token := req.Header.Get("Authorization")
	id, ok = tokens.GetUserId(token)
	return
}

// ParseJson は、リクエストボディをもとに与えられたT型のオブジェクトにJSONをアンマーシャルします。
func ParseJson[T any](req *http.Request, obj *T) (err error) {
	length, err := strconv.Atoi(req.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	jsonBytes := make([]byte, length)
	length, err = req.Body.Read(jsonBytes)
	if err != nil && err != io.EOF {
		return err
	}

	err = json.Unmarshal(jsonBytes[:length], obj)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return
}

// IsNotAuthenticate は、認証されているかを判定し、されていない場合はtrueを返却します。
func IsNotAuthenticate(req *http.Request) bool {
	token := req.Header.Get("Authorization")
	tokens := security.Tokens{}
	_, ok := tokens.GetUserId(token)
	return !ok
}
