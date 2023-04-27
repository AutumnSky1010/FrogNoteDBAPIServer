package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DBConnector は、本番環境のデータベースに接続する構造体です。
type DBConnector struct {
	dbCongig string
}

// Connect は、データベースに接続します。
func (conector *DBConnector) Connect() (database *sql.DB, err error) {
	database, err = sql.Open("mysql", conector.dbCongig)
	if err != nil {
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	return
}

// NewDBConnector は、DBConnector構造体を初期化し、返却します。
func NewDBConnector() *DBConnector {
	pass := os.Getenv("FROGNOTE_DB_PASS")
	// [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
	dbconf := fmt.Sprintf("frognote_db:%s@tcp(localhost:3306)/frognote?charset=utf8mb4", pass)
	return &DBConnector{dbCongig: dbconf}
}
