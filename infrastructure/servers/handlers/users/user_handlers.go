package users

import (
	domainUsers "FrogNote_database/domain/users"
	"FrogNote_database/infrastructure/db"
	dbUsers "FrogNote_database/infrastructure/db/users"
	"FrogNote_database/infrastructure/security"
	"FrogNote_database/infrastructure/servers"
	"FrogNote_database/infrastructure/servers/handlers"
	"encoding/json"
	"net/http"
)

// userObj は、ユーザ情報を表現する構造体です。
type userObj struct {
	Password   string `json:"password"`
	ScreenName string `json:"screenName"`
	SignInId   string `json:"signInId"`
}

// authenticationObj は、認証情報を表現する構造体です。
type authenticationObj struct {
	SignInId string `json:"signInId"`
	Password string `json:"password"`
}

// Modify は、ユーザ情報を編集するためのハンドラです。
func Modify(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if handlers.IsNotJsonReq(req, "PATCH") {
		return http.StatusBadRequest, []byte("Bad request")
	}

	// jsonをパース
	parsedUser := userObj{}
	err := handlers.ParseJson(req, &parsedUser)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not parse json")
	}

	// 既存のユーザ情報を取得]
	repos := dbUsers.NewUserRepository(db.NewDBConnector())
	userId, _ := handlers.GetUserId(req)
	signInId, _ := domainUsers.NewSignInId(parsedUser.SignInId)

	oldUser, err := repos.FindBySignInId(signInId)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not found user")
	}
	// パースしたユーザ情報をもとに組み立てる
	// もし、パスワードが変更されていた場合、ハッシュ化する。
	if parsedUser.Password != oldUser.Password {
		hash := security.Hash{}
		parsedUser.Password = hash.GetHash(parsedUser.Password)
	}
	user, err := domainUsers.NewUser(*userId, parsedUser.ScreenName, *signInId, parsedUser.Password)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not create user")
	}

	// 更新する
	err = repos.Update(user)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not write data")
	}
	return http.StatusOK, []byte("")
}

// Get は、ユーザ情報を取得するためのハンドラです。
func Get(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if req.Method != "GET" {
		return http.StatusBadRequest, []byte("Bad request")
	}

	repos := dbUsers.NewUserRepository(db.NewDBConnector())
	id, _ := handlers.GetUserId(req)
	user, err := repos.FindByUserId(id)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not find user")
	}

	// レスポンス用にオブジェクトを組み立てる。
	resUser := userObj{
		Password:   user.Password,
		ScreenName: user.ScreenName,
		SignInId:   user.SignInId.GetValue(),
	}
	json, err := json.Marshal(resUser)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not convert json.")
	}
	return http.StatusOK, json
}

// Create は、新規ユーザー作成用のハンドラです。
func Create(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotJsonReq(req, "POST") {
		return http.StatusBadRequest, []byte("Bad request")
	}

	repos := dbUsers.NewUserRepository(db.NewDBConnector())
	user := userObj{}
	err := handlers.ParseJson(req, &user)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not parse json")
	}

	signInId, err := domainUsers.NewSignInId(user.SignInId)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte(err.Error())
	}

	_, err = repos.Create(signInId, user.Password, user.ScreenName)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not create user")
	}
	return http.StatusOK, []byte("")
}

// Leave は、ユーザーを削除するためのハンドラです。
func Leave(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if req.Method != "DELETE" {
		return http.StatusBadRequest, []byte("Bad request")
	}

	id, _ := handlers.GetUserId(req)
	repos := dbUsers.NewUserRepository(db.NewDBConnector())
	user, _ := repos.FindByUserId(id)
	err := repos.Delete(&user.SignInId)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not delete user")
	}
	return http.StatusOK, []byte("")
}

// Authenticate は、ユーザ認証を行うためのハンドラです。
func Authenticate(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotJsonReq(req, "POST") {
		return http.StatusBadRequest, []byte("Bad request")
	}
	authObj := authenticationObj{}
	err := handlers.ParseJson(req, &authObj)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not parse json")
	}

	repos := dbUsers.NewUserRepository(db.NewDBConnector())
	signInId, _ := domainUsers.NewSignInId(authObj.SignInId)
	// ユーザ情報を取得
	user, err := repos.FindBySignInId(signInId)

	// リクエスト内のパスワードをハッシュ化
	hash := security.Hash{}
	password := hash.GetHash(authObj.Password)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("internal error")
	}

	// パスワードが違う場合
	if user.Password != password {
		return http.StatusUnauthorized, []byte("Password is incorrect")
	}

	// アクセストークンを生成
	tokens := security.Tokens{}
	token := tokens.GenereteToken(&user.Id)
	return http.StatusOK, []byte(token)
}

// SignOut は、サインアウトするためのハンドラです。セッションを切断します。
func SignOut(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if req.Method != "DELETE" {
		return http.StatusBadRequest, []byte("bad request")
	}
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}

	tokens := security.Tokens{}
	token := req.Header.Get("Authorization")

	// トークンを削除（無効化）する
	tokens.Invalidate(token)
	return http.StatusOK, []byte("")
}

// GetHandlers は、ハンドラのスライスを返却します。
func GetHandlers() []servers.Handler {
	return []servers.Handler{
		{Pattern: "/user/modify", HandlerFunc: Modify},
		{Pattern: "/user/auth", HandlerFunc: Authenticate},
		{Pattern: "/user/create", HandlerFunc: Create},
		{Pattern: "/user/leave", HandlerFunc: Leave},
		{Pattern: "/user/signout", HandlerFunc: SignOut},
		{Pattern: "/user", HandlerFunc: Get},
	}
}
