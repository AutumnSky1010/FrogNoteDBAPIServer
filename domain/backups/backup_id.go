package backups

// backupId は、バックアップIDを表現する構造体です。
type BackupId struct {
	value int
}

// GetValue は、バックアップIDの値を取得します。
func (id *BackupId) GetValue() int {
	return id.value
}

// Equals は、二つのバックアップIDの等価性を判定します。等しければtrueを返却します。
func (id *BackupId) Equals(otherId *BackupId) bool {
	return id.value == otherId.value
}

// NewBackupId は、BackupId構造体を初期化し、返却します。。
func NewBackupId(value int) (id *BackupId) {
	return &BackupId{value: value}
}
