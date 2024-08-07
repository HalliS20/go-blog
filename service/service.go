package service

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type BlogPost struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var Db *sql.DB

func InitDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS blog (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        body TEXT
    );`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBlogPosts() []BlogPost {
	rows, err := Db.Query("SELECT id, title, body FROM blog")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	posts := []BlogPost{}
	for rows.Next() {
		post := BlogPost{}
		err := rows.Scan(&post.ID, &post.Title, &post.Body)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	return posts
}

func CreateBlogPost(post BlogPost) {
	_, err := Db.Exec("INSERT INTO blog (title, body) VALUES (?, ?)", post.Title, post.Body)
	if err != nil {
		log.Fatal(err)
	}
}
