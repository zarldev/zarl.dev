package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/pkg/repo"
	"github.com/zarldev/zarldotdev/view/article"
	"github.com/zarldev/zarldotdev/view/layout"
)

type Article struct {
	Title    string
	Subtitle string
	Body     string
	Image    string
}

type ArticleHandler struct {
	Repo *repo.ArticleRepository
}

func NewArticleHandler(config repo.Config) (*ArticleHandler, error) {
	ar, err := repo.NewArticleRepository(config)
	if err != nil {
		return nil, err
	}
	return &ArticleHandler{
		Repo: ar,
	}, nil
}

func (h *ArticleHandler) RegisterRoutes(router *echo.Echo) {
	router.GET("/articles", h.ListArticles)
	router.GET("/articles/:slug", h.Article)
}

func (h *ArticleHandler) ListArticles(c echo.Context) error {
	articleRows, err := h.Repo.GetPublishedArticles()
	if err != nil {
		return render(c, layout.Error(err))
	}
	articles := make([]article.Article, len(articleRows))
	for pos, a := range articleRows {
		articles[pos] = a.ToArticle()
	}
	return render(c, article.ListArticles(articles))

}

func (h *ArticleHandler) Article(c echo.Context) error {
	slug := c.Param("slug")
	articleRow, err := h.Repo.GetArticleBySlug(slug)
	if err != nil {
		return render(c, layout.Error(err))
	}
	a := articleRow.ToArticle()

	clapped := false
	cookies := c.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == clapID(a.ID) {
			clapped = true
		}
	}
	return render(c, article.Show(a, clapped))
}
