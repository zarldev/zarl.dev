package repo

import (
	"database/sql"
	"fmt"

	svc "github.com/zarldev/zarldotdev/pkg/service"
)

type AdminRepository struct {
	*sql.DB
}

func NewAdminRepository(config Config) (*AdminRepository, error) {
	db, err := sql.Open("sqlite", config.Connection)
	if err != nil {
		return nil, err
	}
	ar := &AdminRepository{db}
	err = ar.createTable()
	if err != nil {
		return nil, err
	}
	if list, err := ar.ListAdmins(); err != nil || len(list) == 0 {
		admin, pass := svc.RandomUsernameAndHashedPassword()
		fmt.Printf("Created admin user: %s with password: %s\n", admin, pass)
	}

	return ar, nil
}

func (a *AdminRepository) createTable() error {
	_, err := a.Exec("CREATE TABLE IF NOT EXISTS admin (username TEXT PRIMARY KEY, password TEXT)")
	if err != nil {
		return err
	}
	return nil
}

func (a *AdminRepository) Get(username string) (string, error) {
	var password string
	err := a.QueryRow("SELECT password FROM admin WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (a *AdminRepository) ListAdmins() ([]string, error) {
	rows, err := a.Query("SELECT username FROM admin")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

func (a *AdminRepository) Set(username, password string) error {
	_, err := a.Exec("INSERT INTO admin (username, password) VALUES (?, ?) ON CONFLICT(username) DO UPDATE SET password = ?", username, password, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdminRepository) Delete(username string) error {
	_, err := a.Exec("DELETE FROM admin WHERE username = ?", username)
	if err != nil {
		return err
	}
	return nil
}
