package snippet

import (
	"database/sql"
	"diskhub/web/logger"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func GetAll() []Snippet {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		logger.Console.Fatal("%e", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content, lang, tags FROM snippets")
	if err != nil {
		logger.Console.Fatal("%e", err)
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		var tags string
		if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Lang, &tags); err != nil {
			logger.Console.Fatal("%e", err)
		}

		s.Loc = len(strings.Split(s.Content, "\n"))
		s.Tags = strings.Split(tags, ",")

		snippets = append(snippets, s)
	}

	return snippets
}

func Get(id string) Snippet {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		logger.Console.Fatal("%e", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content, lang, tags FROM snippets")
	if err != nil {
		logger.Console.Fatal("%e", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Snippet
		var tags string
		if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Lang, &tags); err != nil {
			logger.Console.Fatal("%e", err)
		}

		s.Loc = len(strings.Split(s.Content, "\n"))
		s.Tags = strings.Split(tags, ",")

		if s.Id == id {
			return s
		}
	}

	return Snippet{}
}
