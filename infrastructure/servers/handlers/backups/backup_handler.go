package backups

import (
	domainBackups "FrogNote_database/domain/backups"
	"FrogNote_database/infrastructure/db"
	dbBackups "FrogNote_database/infrastructure/db/backups"
	"FrogNote_database/infrastructure/servers"
	"FrogNote_database/infrastructure/servers/handlers"
	"encoding/json"
	"io"
	"net/http"
)

// backupId は、バックアップIDを表現する構造体です。
type backupId struct {
	Value int `json:"value"`
}

// Delete は、バックアップデータ（本体・メタデータ）を削除するハンドラです。
func Delete(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if handlers.IsNotJsonReq(req, "DELETE") {
		return http.StatusBadRequest, []byte("Bad request")
	}

	parsedId, err := parseBackupId(req)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not parse json")
	}

	repos := dbBackups.NewBackupRepository(db.NewDBConnector())
	backup, err := repos.FindByBackupId(parsedId)
	// バックアップデータが見つからなかった場合。
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Not found backup")
	}

	// もしバックアップが別ユーザのものだった場合は認証されていないという扱いとする。
	if !isOwner(req, backup) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}

	// バックアップデータを削除する。
	err = repos.DeleteByBackupId(parsedId)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not delete.")
	}
	return http.StatusOK, []byte("")
}

// Save は、バックアップデータ（本体のバイナリ）を保存するハンドラです。
func Save(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if req.Method != "POST" {
		return http.StatusBadRequest, []byte("Bad request")
	}

	userId, _ := handlers.GetUserId(req)
	repos := dbBackups.NewBackupRepository(db.NewDBConnector())

	file, _, err := req.FormFile("backup")
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not get backup file.")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not read backup file.")
	}

	err = repos.Create(userId, data)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not save backup.")
	}
	return http.StatusOK, []byte("")
}

// Download は、バックアップデータ（本体のバイナリ）をダウンロードするためのハンドラです。
func Download(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	if req.Method != "POST" {
		return http.StatusBadRequest, []byte("Bad request")
	}

	id, err := parseBackupId(req)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not parse json")
	}

	repos := dbBackups.NewBackupRepository(db.NewDBConnector())
	backup, err := repos.FindByBackupId(id)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Not found backup")
	}

	// もしバックアップが別ユーザのものだった場合は認証されていないという扱いとする。
	if !isOwner(req, backup) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	return http.StatusOK, backup.Backup
}

// GetAllmeta は、ユーザが所有するすべてのバックアップデータのメタデータを取得するためのハンドラです。
func GetAllmeta(writer http.ResponseWriter, req *http.Request, logger *servers.Logger) (status int, body []byte) {
	if handlers.IsNotAuthenticate(req) {
		return http.StatusUnauthorized, []byte("Unauthorized")
	}
	userid, _ := handlers.GetUserId(req)
	repos := dbBackups.NewBackupRepository(db.NewDBConnector())
	backups, err := repos.FindBackupMetas(userid)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not find backups")
	}

	// backupMeta は、レスポンス用のバックアップメタデータを表現する構造体
	type backupMeta struct {
		BackupId int    `json:"backupId"`
		SavedAt  string `json:"savedAt"`
	}

	// レスポンス用のメタデータ構造体に詰め替える
	metas := make([]backupMeta, len(backups))
	for i, backup := range backups {
		metas[i] = backupMeta{BackupId: backup.BackupId.GetValue(), SavedAt: backup.SavedAt}
	}
	json, err := json.Marshal(metas)
	if err != nil {
		logger.FPrintErrorLog(err, "")
		return http.StatusInternalServerError, []byte("Could not convert to json")
	}
	return http.StatusOK, json
}

// isOwner は、引数のバックアップデータをリクエスト送信元ユーザ所有かを判定し、所有している場合はtrueを返します。
func isOwner(req *http.Request, backup *domainBackups.Backup) bool {
	userid, _ := handlers.GetUserId(req)
	return backup.UserId.Equals(userid)
}

// parseBackupId は、リクエストボディからバックアップIDをパースして返します。
func parseBackupId(req *http.Request) (id *domainBackups.BackupId, err error) {
	parsedId := backupId{}
	err = handlers.ParseJson(req, &parsedId)
	if err != nil {
		return nil, err
	}
	id = domainBackups.NewBackupId(parsedId.Value)
	return id, nil
}

// GetHandlers は、ハンドラのスライスを返却します。
func GetHandlers() []servers.Handler {
	return []servers.Handler{
		{Pattern: "/backup/save", HandlerFunc: Save},
		{Pattern: "/backup/delete", HandlerFunc: Delete},
		{Pattern: "/backup/allmeta", HandlerFunc: GetAllmeta},
		{Pattern: "/backup/download", HandlerFunc: Download},
	}
}
