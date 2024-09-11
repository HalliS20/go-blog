package interfaces

import "go-blog/internal/domain/models"

type BlogPost = models.BlogPost

type BlogPostServiceInterface interface {
	GetBlogPosts() []BlogPost
	CreateBlogPost(post BlogPost)
}
