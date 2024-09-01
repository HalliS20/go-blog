package models

type BlogPost struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Date        string `json:"date"`
}

func (BlogPost) TableName() string {
	return "posts"
}

func MakeBlogPost(title string, description string, body string) BlogPost {
	return BlogPost{
		Title:       title,
		Description: description,
		Body:        body,
	}
}
