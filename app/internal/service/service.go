// service/service.go
package service

import (
	"database/sql"
	"go-blog/internal/models"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var Db *sql.DB

type BlogPost = models.BlogPost

func InitDatabase() {
	user := "Blog_owner"
	password := os.Getenv("DB_PASSWORD")
	connStr := "postgresql://" + user + ":" + password + "@ep-late-sun-a5p8yfr7.us-east-2.aws.neon.tech/Blog?sslmode=require"
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal("Error connecting to the database:", err)
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

		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Body, &date)
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
	_, err := Db.Exec("INSERT INTO posts (title, description, body) VALUES ($1, $2, $3)", post.Title, post.Description, post.Body)
	if err != nil {
		log.Fatal("Error inserting post: ", err)
	}
}

func CloseDatabase() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Println("Error closing database connection:", err)
		}
	}
}
