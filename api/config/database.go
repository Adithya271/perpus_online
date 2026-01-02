package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/perpustakaan_digital")
	if err != nil {
		panic(err)
	}
	return db
}
