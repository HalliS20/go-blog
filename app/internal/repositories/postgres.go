package repositories

import (
	"go-blog/internal/domain/interfaces"
	"go-blog/internal/domain/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type (
	BlogPost           = models.BlogPost
	BlogPostRepository = interfaces.BlogPostRepository
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(DB *gorm.DB) BlogPostRepository {
	return &PostgresRepository{DB: DB}
}

func (r *PostgresRepository) GetBlogPosts() []BlogPost {
	var posts []BlogPost
	result := r.DB.Order("id desc").Find(&posts)
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

func (r *PostgresRepository) GetBlogPost(id int) BlogPost {
	post := BlogPost{}
	result := r.DB.First(&post, id)
	if result.Error != nil {
		log.Fatal("Error querying post: ", result.Error)
	}
	return post
}

func (r *PostgresRepository) CreateBlogPost(post BlogPost) error {
	if post.Date == "" {
		post.Date = time.Now().Format("2006-01-02") // Set the current date if not provided
	}
	result := r.DB.Create(&post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
