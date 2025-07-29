package migrations

import (
	"database/sql"
	"log"

	"github.com/mattn/go-sqlite3"
)

func InitDb(db *sql.DB) error {
	var version uint
	for {
		err := db.QueryRow(`SELECT value FROM migrations WHERE key = 'version'`).Scan(&version)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok {
				switch err.Code {
				case sqlite3.ErrError:
					version = 0
				default:
					return err
				}
			} else {
				return err
			}
		}

		switch version {
		case 0:
			log.Println("Initializing database")
			err := MigrationInit(db)
			if err != nil {
				return err
			}
		case 1:
			return nil
		default:
			log.Fatalf("Unknown database schema version: %d. Cannot continue.", version)
		}
	}
}
