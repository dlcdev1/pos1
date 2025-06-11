package database

import (
	"database/sql"
)

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	)`)
	return err
}

//func InitDBCourses(db *sql.DB) error {
//	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS courses (
//		id TEXT PRIMARY KEY,
//		name TEXT NOT NULL,
//		description TEXT,
//		category_id TEXT NOT NULL,
//		FOREIGN KEY (category_id) REFERENCES categories (id)
//	)`)
//	return err
//}
