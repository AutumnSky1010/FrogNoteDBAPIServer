package users

import (
	"FrogNote_database/domain/users"
	"FrogNote_database/infrastructure/db"
	"FrogNote_database/infrastructure/security"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

// UserRepository は、ユーザを永続化・復元する構造体です。
type UserRepository struct {
	connector db.IDBConnector
}

// Create はユーザを新規保存します。
func (repos *UserRepository) Create(signInId *users.SignInId, password string, screenName string) (user *users.User, err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// パスワードをハッシュ化する。
	hash := security.Hash{}
	password = hash.GetHash(password)
	_, err = db.Exec("insert into users values (default, ?, ?, ?)", signInId.GetValue(), password, screenName)
	if err != nil {
		return nil, err
	}
	user, err = repos.FindBySignInId(signInId)
	return
}

// Update は、ユーザ情報を更新します。
func (repos *UserRepository) Update(user *users.User) (err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	old, err := repos.FindByUserId(&user.Id)
	if err != nil {
		return errors.New("not found user")
	}

	password := user.Password
	// パスワードが変更された場合、ハッシュ値にする。
	if old.Password != password {
		hash := security.Hash{}
		password = hash.GetHash(password)
	}

	_, err = db.Exec("update users set sign_in_id = ?, password = ?, screen_name = ? where id = ?", user.SignInId.GetValue(), password, user.ScreenName, user.Id.GetValue())
	return
}

// FindBySignInId は、サインインIDをもとにユーザを取得します。
func (repos *UserRepository) FindBySignInId(signInId *users.SignInId) (user *users.User, err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("select * from users where users.sign_in_id = ?", signInId.GetValue())
	user, err = mapUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByUserId は、ユーザIDをもとにユーザを取得します。
func (repos *UserRepository) FindByUserId(userId *users.UserId) (user *users.User, err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	row := db.QueryRow("select * from users where users.id = ?", userId.GetValue())
	user, err = mapUser(row)
	if err != nil {
		return nil, err
	}
	return
}

// Delete は、指定したサインインIDのユーザを削除します。1
func (repos *UserRepository) Delete(signInId *users.SignInId) (err error) {
	db, err := repos.connector.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from users where users.sign_in_id = ?", signInId.GetValue())
	if err != nil {
		return err
	}
	return nil
}

// mapUser は、rowから読み取ります。
func mapUser(row *sql.Row) (user *users.User, err error) {
	var userIdValue int
	var signInIdStr string
	user = &users.User{}

	err = row.Scan(&userIdValue, &signInIdStr, &user.Password, &user.ScreenName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	readUserId := users.NewUserId(userIdValue)
	readSignInId, _ := users.NewSignInId(signInIdStr)
	user.Id = *readUserId
	user.SignInId = *readSignInId
	return user, nil
}

// NewUserRepository は、UserRepositoryを初期化します。
func NewUserRepository(connector db.IDBConnector) (repos *UserRepository) {
	return &UserRepository{connector: connector}
}
