package db

import "database/sql"

// IDBConnector は、DBに接続するコネクタのインターフェースです。
type IDBConnector interface {
	// Connect は、データベースに接続します。
	Connect() (database *sql.DB, err error)
}
