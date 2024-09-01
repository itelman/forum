package store

import (
	"database/sql"
	"io/ioutil"
)

func Migrate(db *sql.DB, migrPath string) error {
	content, err := ioutil.ReadFile(migrPath)
	if err != nil {
		return err
	}
	script := string(content)

	_, err = db.Exec(script)
	if err != nil {
		return err
	}

	return nil
}
