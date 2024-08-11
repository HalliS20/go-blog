// service/service.go
package service

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type BlogPost struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

var Db *sql.DB

func InitDatabase() {
	connStr := "postgresql://Blog_owner:D4nb2hMustHr@ep-late-sun-a5p8yfr7.us-east-2.aws.neon.tech/Blog?sslmode=require"
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBlogPosts() []BlogPost {
	rows, err := Db.Query("SELECT id, title, body, description FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		post := BlogPost{}
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Description)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	return posts
}

func CreateBlogPost(post BlogPost) {
}
