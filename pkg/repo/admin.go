package repo

import (
	"fmt"

	svc "github.com/zarldev/zarldotdev/pkg/service"
)

type AdminsRepository interface {
	Get(username string) (string, error)
	ListAdmins() ([]string, error)
	Set(username, password string) error
	Delete(username string) error
}

type SQLiteAdminRepo struct {
	conn *Connection
}

// tar extract
func NewAdminRepository(conn *Connection) (*SQLiteAdminRepo, error) {
	ar := &SQLiteAdminRepo{conn: conn}
	err := ar.createTable()
	if err != nil {
		return nil, err
	}
	if list, err := ar.ListAdmins(); err != nil || len(list) == 0 {
		admin, pass, cryptPass := svc.AdminPassCrypted()
		fmt.Printf("Created admin user: %s with password: %s - bcrypted to %s \n", admin, pass, cryptPass)
		err := ar.Set(admin, cryptPass)
		if err != nil {
			return nil, err
		}
	}
	return ar, nil
}

func (a *SQLiteAdminRepo) createTable() error {
	d, err := a.conn.write.Exec("CREATE TABLE IF NOT EXISTS admin (username TEXT PRIMARY KEY, password TEXT)")
	fmt.Println(d)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdminRepo) Get(username string) (string, error) {
	var password string
	err := a.conn.read.QueryRow("SELECT password FROM admin WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (a *SQLiteAdminRepo) ListAdmins() ([]string, error) {
	rows, err := a.conn.read.Query("SELECT username FROM admin")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	var admins []string
	for rows.Next() {
		var admin string
		err := rows.Scan(&admin)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

func (a *SQLiteAdminRepo) Set(username, password string) error {
	_, err := a.conn.write.Exec("INSERT INTO admin (username, password) VALUES (?, ?) ON CONFLICT(username) DO UPDATE SET password = ?", username, password, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdminRepo) Delete(username string) error {
	_, err := a.conn.write.Exec("DELETE FROM admin WHERE username = ?", username)
	if err != nil {
		return err
	}
	return nil
}
