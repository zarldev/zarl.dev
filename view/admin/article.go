package admin

import "time"

type Article struct {
	ID        int
	Title     string
	Subtitle  string
	Content   string
	Image     string
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published bool
	Markdown  string
	Claps     int
}

type CreateArticle struct {
	Title     string `form:"title"`
	Subtitle  string `form:"subtitle"`
	Image     string `form:"image"`
	Slug      string `form:"slug"`
	Markdown  string `form:"markdown"`
	Published bool   `form:"published"`
}
