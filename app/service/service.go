// service/service.go
package service

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type BlogPost struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Date        string `json:"date"`
}

var Db *sql.DB

func InitDatabase() {
	connStr := "postgresql://Blog_owner:D4nb2hMustHr@ep-late-sun-a5p8yfr7.us-east-2.aws.neon.tech/Blog?sslmode=require"
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        body TEXT NOT NULL,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBlogPosts() []BlogPost {
	rows, err := Db.Query("SELECT * FROM posts ORDER BY id DESC")
	if err != nil {
		log.Fatal("Error querying posts table: ", err)
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		post := BlogPost{}
		var date time.Time

		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Description, &date)
		if err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		post.Date = date.Format("02 / 01 / 2006")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error iterating over rows: ", err)
	}

	return posts
}


func CreateBlogPost(post BlogPost) {
}
