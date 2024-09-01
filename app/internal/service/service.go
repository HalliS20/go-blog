// service/service.go
package service

import (
	"go-blog/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

type BlogPost = models.BlogPost

var (
	user     string
	password string
	connStr  string
)

func InitDatabase() {
	// set variables
	user = "Blog_owner"
	password = os.Getenv("DB_PASSWORD")
	connStr = "host=ep-late-sun-a5p8yfr7.us-east-2.aws.neon.tech user=" + user + " password=" + password + " dbname=Blog port=5432 sslmode=require TimeZone=UTC"

	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = DB.AutoMigrate(&BlogPost{})
	if err != nil {
		log.Fatal("Error migrating database schema:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error getting underlying SQL DB:", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

}

func GetBlogPosts() []BlogPost {
	var posts []BlogPost
	result := DB.Order("id desc").Find(&posts)
	if result.Error != nil {
		log.Fatal("Error querying posts: ", result.Error)
	}

	for i := range posts {
		// Check if the Date field has content before slicing
		if len(posts[i].Date) >= 10 {
			posts[i].Date = posts[i].Date[:10] // Extract only the date part (YYYY-MM-DD)
		} else {
			// If the Date field is empty or too short, set it to the current date
			posts[i].Date = time.Now().Format("2006-01-02")
			log.Printf("Warning: Post ID %d had an invalid date. Set to current date.\n", posts[i].ID)
		}
	}

	return posts
}

func CreateBlogPost(post BlogPost) {
	if post.Date == "" {
		post.Date = time.Now().Format("2006-01-02") // Set the current date if not provided
	}
	result := DB.Create(&post)
	if result.Error != nil {
		log.Fatal("Error inserting post: ", result.Error)
	}
}

func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting underlying SQL DB:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Error closing database connection:", err)
	}

	log.Println("Database and listener closed successfully")
}
