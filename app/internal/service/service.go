// service/service.go
package service

import (
	"go-blog/internal/domain/interfaces"
	"go-blog/internal/domain/models"
	"log"
)

type (
	BlogPost              = models.BlogPost
	BlogPostRepoInterface = interfaces.BlogPostRepoInterface
)

type BlogService struct {
	repo BlogPostRepoInterface
}

func NewBlogService(repo BlogPostRepoInterface) interfaces.BlogPostServiceInterface {
	return &BlogService{repo: repo}
}

func (s *BlogService) GetBlogPosts() []BlogPost {
	posts := s.repo.GetBlogPosts()
	return posts
}

func (s *BlogService) CreateBlogPost(post BlogPost) {
	err := s.repo.CreateBlogPost(post)
	if err != nil {
		log.Println("Error creating blog post:", err)
	}
}
