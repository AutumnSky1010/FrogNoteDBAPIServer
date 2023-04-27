package backups

import (
	"FrogNote_database/domain/backups"
	"FrogNote_database/domain/users"
	"FrogNote_database/infrastructure/db"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// BackupRepository は、バックアップを永続化・復元する構造体です。
type BackupRepository struct {
	connector db.IDBConnector
}

// FindBackupMetas は、バックアップのメタデータのスライスを取得します。（バックアップの本体がない状態で返却されます。）
func (repos *BackupRepository) FindBackupMetas(userId *users.UserId) (backupSlice []*backups.Backup, err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, user_id, saved_at from backups where backups.user_id = ?", userId.GetValue())
	if err != nil {
		return nil, err
	}
	backupSlice, err = mapBackups(rows, true)
	if err != nil {
		return nil, err
	}
	return backupSlice, nil
}

// FindByBackupId は、BackupIdをもとにバックアップを取得します。
func (repos *BackupRepository) FindByBackupId(backupId *backups.BackupId) (backup *backups.Backup, err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select * from backups where backups.id = ?", backupId.GetValue())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	backupSlice, err := mapBackups(rows, false)
	if err != nil {
		return nil, err
	}
	return backupSlice[0], nil
}

// DeleteByBackupId は、バックアップIDをもとにバックアップを削除します。
func (repos *BackupRepository) DeleteByBackupId(backupId *backups.BackupId) (err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from backups where backups.id = ?", backupId.GetValue())
	return err
}

// Create は、バックアップを新規保存します。
func (repos *BackupRepository) Create(userId *users.UserId, backupBin []byte) error {

	db, err := repos.connector.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("insert into backups values (default, ?, ?, default)", userId.GetValue(), backupBin)
	return err
}

// mapBackupsは、複数のバックアップデータをマップして返却します。
func mapBackups(rows *sql.Rows, isMetaOnly bool) (backupSlice []*backups.Backup, err error) {
	backupSlice = make([]*backups.Backup, 0)
	for rows.Next() {
		backup := &backups.Backup{}
		var userIdValue int
		var backupIdValue int
		if isMetaOnly {
			rows.Scan(&backupIdValue, &userIdValue, &backup.SavedAt)
		} else {
			rows.Scan(&backupIdValue, &userIdValue, &backup.Backup, &backup.SavedAt)
		}
		backupId := backups.NewBackupId(backupIdValue)
		userId := users.NewUserId(userIdValue)
		backup.UserId = *userId
		backup.BackupId = *backupId
		backupSlice = append(backupSlice, backup)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return backupSlice, nil
}

// NewBackupRepository は、BackupRepository構造体を初期化し、返却します。
func NewBackupRepository(connector db.IDBConnector) (repos *BackupRepository) {
	return &BackupRepository{connector: connector}
}
