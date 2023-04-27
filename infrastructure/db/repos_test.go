// db_test は、リポジトリのテスト用パッケージです。単体テストではなく、結合テストにあたります。また、テストを行う順番には制約があります。
package db_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	dom_backups "FrogNote_database/domain/backups"
	dom_users "FrogNote_database/domain/users"
	inf_backups "FrogNote_database/infrastructure/db/backups"
	inf_users "FrogNote_database/infrastructure/db/users"
	"FrogNote_database/infrastructure/security"

	_ "github.com/go-sql-driver/mysql"
)

// TestDBConnector は、テスト用データベースに接続するための構造体です。
type TestDBConnector struct {
	dbCongig string
}

// Connect は、テスト用データベースに接続します。
func (conector *TestDBConnector) Connect() (database *sql.DB, err error) {
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

// NewTestDBConnector は、TestDBConnector構造体を初期化し、返却します。
func NewTestDBConnector() *TestDBConnector {
	pass := os.Getenv("FROGNOTE_DB_TESTER_PASS")
	// [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
	dbconf := fmt.Sprintf("frognote_db_tester:%s@tcp(localhost:3306)/frognote_test?charset=utf8mb4", pass)
	return &TestDBConnector{dbCongig: dbconf}
}

var (
	// dummyUser1Id は、あらかじめusersテーブルに追加したダミーユーザのIDです。
	dummyUser1Id = dom_users.NewUserId(1)
	// dummyUser1 は、あらかじめusersテーブルに追加したダミーユーザです。
	dummyUser1 = dom_users.User{
		Id:         *dummyUser1Id,
		SignInId:   *getDummyUser1SignInId(),
		ScreenName: "dummy_1",
		Password:   "255",
	}

	// dummyBackup1 は、あらかじめbackupsテーブルに追加したダミーバックアップです。
	dummyBackup1 = dom_backups.Backup{
		BackupId: *dom_backups.NewBackupId(1),
		UserId:   *dummyUser1Id,
		Backup:   []byte{255},
		SavedAt:  "2023-04-09 13:51:13",
	}

	backupRepos = inf_backups.NewBackupRepository(NewTestDBConnector())
	userRepos   = inf_users.NewUserRepository(NewTestDBConnector())
)

func getDummyUser1SignInId() *dom_users.SignInId {
	signInId, _ := dom_users.NewSignInId("11111111")
	return signInId
}

func getDummyUser2SignInId() *dom_users.SignInId {
	dummy2Signinid, _ := dom_users.NewSignInId("12345678")
	return dummy2Signinid
}

// TestStart は、テストのエントリポイントです。
func TestStart(t *testing.T) {
	/**
	 * テストを行う順番には制約があります。
	 * 1．ユーザを探し出せるかを証明します
	 * 2. バックアップを探し出せるかを証明します
	 * 3. ユーザのライフサイクルをもとにテストします
	 * 3.1 ユーザが作成できるかをテストします
	 * 3.2 ユーザが更新されるかをテストします
	 * 3.3 バックアップを保存できるかをテストします
	 * 3.4 バックアップを削除できるかをテストします
	 * 3.5 ユーザを削除できるかをテストします
	 */
	t.Run("サインインIDをもとにユーザを探す", testFindUserBySignInId)
	t.Run("ユーザIDをもとにユーザを探す", testFindUserByUserId)
	t.Run("ユーザIDをもとにバックアップを探す", testFindBackupByUserId)
	t.Run("バックアップIDをもとにバックアップを探す。", testFindBackupByBackupId)
	t.Run("ユーザが作成できるか", testCreateUser)
}

// testFindUserByUserId は、ユーザIDをもとにユーザを取得できるかをテストします。
func testFindUserByUserId(t *testing.T) {
	dummyUser1FromRepos, err := userRepos.FindByUserId(&dummyBackup1.UserId)
	if err != nil {
		t.Error(err)
	}
	if equalsUser(dummyUser1FromRepos, &dummyUser1) {
		t.Log("pass")
	} else {
		t.Error("invalid user data.")
	}
}

// testFindBackupByBackupId は、バックアップIDをもとに、バックアップを取得できるかをテストします。
func testFindBackupByBackupId(t *testing.T) {
	dummyFromRepos, err := backupRepos.FindByBackupId(&dummyBackup1.BackupId)
	if err != nil {
		t.Error(err)
		return
	}
	if equalsBackup(dummyFromRepos, &dummyBackup1) {
		t.Log("pass")
	} else {
		t.Error("invalid backup data.")
	}
}

// testFindUserBySignInId は、サインインIDをもとにユーザを取得できるかをテストします。
func testFindUserBySignInId(t *testing.T) {
	dummyUser1FromRepos, err := userRepos.FindBySignInId(&dummyUser1.SignInId)
	if err != nil {
		t.Error(err)
	}
	if equalsUser(dummyUser1FromRepos, &dummyUser1) {
		t.Log("pass")
	} else {
		t.Error("invalid user data.")
	}
}

// testFindBackupByUserId は、ユーザIdをもとにバックアップデータを取得できるかをテストします。
func testFindBackupByUserId(t *testing.T) {
	dummiesFromRepos, err := backupRepos.FindBackupMetas(dummyUser1Id)
	if err != nil {
		t.Error(err)
		return
	}
	dummyFromRepos := dummiesFromRepos[0]
	// メタデータが取得されるため、バックアップデータ本体は入っていないので追加する。
	dummyFromRepos.Backup = dummyBackup1.Backup
	if equalsBackup(dummyFromRepos, &dummyBackup1) {
		t.Log("pass")
	} else {
		t.Error("invalid backup data.")
	}
}

// testCreateUser は、ユーザを作成できるかをテストします。
func testCreateUser(t *testing.T) {
	password := "password"
	screenName := "screenName"
	user, err := userRepos.Create(getDummyUser2SignInId(), password, screenName)
	if err != nil {
		t.Error(err)
		return
	}
	// パスワードはハッシュ化される
	hash := security.Hash{}
	password = hash.GetHash(password)

	if user.Password == password && user.ScreenName == screenName && user.SignInId.Equals(getDummyUser2SignInId()) {
		t.Log("pass")
		t.Run("ユーザが更新されるか", testUpdateUser)
	}
}

// testUpdateUser は、ユーザを更新できるかをテストします。
func testUpdateUser(t *testing.T) {
	screenName := "updatedScreenName"
	dummy2, _ := userRepos.FindBySignInId(getDummyUser2SignInId())
	dummy2.ScreenName = screenName
	err := userRepos.Update(dummy2)
	if err != nil {
		t.Error(err)
		return
	}
	newDummy2, _ := userRepos.FindBySignInId(getDummyUser2SignInId())
	if newDummy2.ScreenName != dummy2.ScreenName {
		t.Error("could not update.")
		return
	}
	t.Log("pass")
	t.Run("バックアップできるか", testCreateBackup)
}

// testCreateBackup は、バックアップを作成できるかをテストします。
func testCreateBackup(t *testing.T) {
	dummyUser2, _ := userRepos.FindBySignInId(getDummyUser2SignInId())
	backup := []byte{255}
	backupRepos.Create(&dummyUser2.Id, backup)

	dummies, err := backupRepos.FindBackupMetas(&dummyUser2.Id)
	if err != nil {
		t.Errorf("could not find dummyBackup2.\n %s", err.Error())
	}
	// 空じゃないか
	if len := len(dummies); len != 0 {
		dummy2 := dummies[len-1]
		t.Log("pass")
		t.Run("バックアップを削除できるか", func(t *testing.T) {
			testDeleteBackupByBackupId(t, dummy2.BackupId)
		})
		return
	}
	t.Errorf("could not create dummy2.")
}

// testDeleteBackupByBackupId は、バックアップIDをもとにバックアップを削除できるかをテストします。
func testDeleteBackupByBackupId(t *testing.T, dummy2Id dom_backups.BackupId) {
	dummyUser2, _ := userRepos.FindBySignInId(getDummyUser2SignInId())

	old, _ := backupRepos.FindBackupMetas(&dummyUser2.Id)
	oldLength := len(old)
	err := backupRepos.DeleteByBackupId(&dummy2Id)
	if err != nil {
		t.Errorf("could not delete. please delete id=%d. %s", dummy2Id.GetValue(), err.Error())
	}
	new, _ := backupRepos.FindBackupMetas(&dummyUser2.Id)
	newLength := len(new)
	if oldLength-newLength != 1 {
		t.Errorf("could not delete. please delete id=%d. %s", dummy2Id.GetValue(), err.Error())
	} else {
		t.Log("pass")
		t.Run("ユーザを削除できるか", testDeleteUser)
	}
}

// testDeleteUser は、ユーザを削除できるかをテストします。また、削除されたユーザのバックアップデータがなくなっているかもテストします。
func testDeleteUser(t *testing.T) {
	dummyUser2, _ := userRepos.FindBySignInId(getDummyUser2SignInId())
	// ダミーバックアップデータを作成する
	backupRepos.Create(&dummyUser2.Id, []byte{255})
	backupRepos.Create(&dummyUser2.Id, []byte{255})
	backupRepos.Create(&dummyUser2.Id, []byte{255})

	err := userRepos.Delete(&dummyUser2.SignInId)

	if err != nil {
		t.Errorf("could not delete. please delete id=%d. %s", dummyUser2.Id.GetValue(), err.Error())
	}

	backupSlice, err := backupRepos.FindBackupMetas(&dummyUser2.Id)
	if err != nil {
		t.Log("pass")
	}
	if len := len(backupSlice); len != 0 {
		t.Errorf("could not delete backups of dummyUser2. dummyUser2.Id = %d, %s", dummyUser2.Id.GetValue(), err.Error())
	}
	t.Log("pass")
}

// equalsBackup は、バックアップの全フィールドをもとにを等価比較します。
func equalsBackup(x *dom_backups.Backup, y *dom_backups.Backup) bool {
	return x.UserId.Equals(&y.UserId) &&
		x.BackupId.Equals(&y.BackupId) &&
		x.SavedAt == y.SavedAt &&
		x.Backup[0] == y.Backup[0]
}

// equalsUser は、ユーザの全フィールドをもとにを等価比較します。
func equalsUser(x *dom_users.User, y *dom_users.User) bool {
	return x.Id.Equals(&y.Id) &&
		x.SignInId.Equals(&y.SignInId) &&
		x.Password == y.Password &&
		x.ScreenName == y.ScreenName
}
