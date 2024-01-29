package article

import "time"

type Article struct {
	ID       int
	Title    string
	Image    string
	Subtitle string
	Content  string
	Slug     string
	Created  time.Time
	Claps    int
}
