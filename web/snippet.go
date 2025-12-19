package main

import (
	"strings"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Snippet struct {
	Title 	string
	Content string
	Lang 	string
	Loc 	int
	Tags 	[]string
	Id      string
}

func InitSnippetDb() {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		console.fatal("%e", err)
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
		console.fatal("%e", err)
	}

	db.Close()
}

func GetAllSnippet() []Snippet {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		console.fatal("%e", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content, lang, tags FROM snippets")
	if err != nil { console.fatal("%e", err) }
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
        var s Snippet
		var tags string
        if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Lang, &tags); err != nil {
            console.fatal("%e", err)
        }

		s.Loc = len(strings.Split(s.Content, "\n"))
		s.Tags = strings.Split(tags, ",")

		snippets = append(snippets, s)
    }

	return snippets
}

func GetSnippet(id string) Snippet {
	db, err := sql.Open("sqlite3", "./snippet.db")
	if err != nil {
		console.fatal("%e", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content, lang, tags FROM snippets")
	if err != nil { console.fatal("%e", err) }
	defer rows.Close()

	for rows.Next() {
        var s Snippet
		var tags string
        if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Lang, &tags); err != nil {
            console.fatal("%e", err)
        }

		s.Loc = len(strings.Split(s.Content, "\n"))
		s.Tags = strings.Split(tags, ",")

		if s.Id == id { return s }
    }

	return Snippet{}
}