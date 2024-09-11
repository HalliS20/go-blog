// service/service.go
package service

import (
	"go-blog/internal/domain/interfaces"
	"go-blog/internal/domain/models"
	"log"
)

type (
	BlogPost           = models.BlogPost
	BlogPostRepository = interfaces.BlogPostRepository
)

type BlogService struct {
	repo BlogPostRepository
}

func NewBlogService(repo BlogPostRepository) *BlogService {
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
