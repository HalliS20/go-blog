package interfaces

import "go-blog/internal/domain/models"

type BlogPost = models.BlogPost

type BlogService interface {
	GetBlogPosts() []BlogPost
	CreateBlogPost(post BlogPost)
}
