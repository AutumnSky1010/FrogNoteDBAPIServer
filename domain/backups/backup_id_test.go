package backups_test

import (
	"FrogNote_database/domain/backups"
	"testing"
)

func TestGetValue(t *testing.T) {
	value := 1
	id := backups.NewBackupId(value)
	if id.GetValue() != value {
		t.Error()
	}
}

func TestEquals(t *testing.T) {
	value := 1
	id1 := backups.NewBackupId(value)
	id2 := backups.NewBackupId(value)

	if !id1.Equals(id2) {
		t.Error()
	}
}
