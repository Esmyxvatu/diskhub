package snippet

import (
	"database/sql"
	"diskhub/web/logger"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		logger.Console.Fatal("%e", err)
	}
	defer db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS snippets (
			id TEXT PRIMARY KEY,
			title TEXT,
			content TEXT,
			lang TEXT,
			tags TEXT
		);
	`
	if _, err := db.Exec(createTable); err != nil {
		logger.Console.Fatal("%e", err)
	}

	db.Close()
}
