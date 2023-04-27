package backups

import "FrogNote_database/domain/users"

// Backup は、バックアップデータを表現する構造体です。
type Backup struct {
	BackupId BackupId
	UserId   users.UserId
	Backup   []byte
	SavedAt  string
}

// NewBackup は、バックアップ構造体を初期化し、返却します。
func NewBackup(backupId BackupId, userId users.UserId, savedAt string, backupBlob []byte) (backup *Backup) {
	return &Backup{
		BackupId: backupId,
		UserId:   userId,
		Backup:   backupBlob,
		SavedAt:  savedAt,
	}
}
