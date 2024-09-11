package interfaces

import "go-blog/internal/domain/models"

type BlogPostRepoInterface interface {
	GetBlogPosts() []models.BlogPost
	GetBlogPost(id int) models.BlogPost
	CreateBlogPost(post models.BlogPost) error
}
