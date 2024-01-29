package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/pkg/repo"
	"github.com/zarldev/zarldotdev/view/admin"
	"github.com/zarldev/zarldotdev/view/article"
	"github.com/zarldev/zarldotdev/view/layout"
)

type AdminHandler struct {
	ArticleRepo *repo.ArticleRepository
	AdminRepo   *repo.AdminRepository
}

func RowToArticle(a repo.ArticleRow) article.Article {
	return article.Article{
		Title:    a.Title,
		Subtitle: a.Subtitle,
		Content:  a.Body,
		Image:    a.Image,
		Slug:     a.Slug,
	}
}

var ErrInitialisingHandler = fmt.Errorf("error initialising admin handler")

func NewAdminHandler(config repo.Config) (*AdminHandler, error) {
	ar, err := repo.NewArticleRepository(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInitialisingHandler, err)
	}
	adr, err := repo.NewAdminRepository(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInitialisingHandler, err)
	}
	return &AdminHandler{
		ArticleRepo: ar,
		AdminRepo:   adr,
	}, nil
}

func (h *AdminHandler) RegisterRoutes(router *echo.Echo) {
	admingroup := router.Group("/admin")
	admingroup.GET("", h.AdminLoginShow)
	admingroup.POST("/login", h.AdminLogin)
	protected := admingroup.Group("")
	protected.Use(h.AdminAuth)
	protected.GET("/articles", h.AdminArticlesShow)
	protected.GET("/articles/:id", h.AdminEditArticle)
	protected.POST("/articles/:id/save", h.AdminEditArticleSave)
	protected.GET("/articles/new", h.AdminNewArticle)
	protected.POST("/articles/new", h.AdminNewArticleCreate)
}

func (h *AdminHandler) AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("admin")
		if err != nil {
			return c.Redirect(302, "/admin")
		}
		_, err = h.AdminRepo.Get(cookie.Value)
		if err != nil {
			return c.Redirect(302, "/admin")
		}
		return next(c)
	}
}

func (h *AdminHandler) AdminNewArticleCreate(c echo.Context) error {
	var article admin.CreateArticle
	if err := c.Bind(&article); err != nil {
		return render(c, layout.Error(err))
	}
	markdown := article.Markdown
	body := mdToHTML([]byte(markdown))

	articleRow := repo.ArticleRow{
		Slug:         article.Slug,
		Title:        article.Title,
		Subtitle:     article.Subtitle,
		Body:         string(body),
		MarkdownBody: article.Markdown,
		Image:        article.Image,
		Published:    true,
	}
	err := h.ArticleRepo.CreateArticle(&articleRow)
	if err != nil {
		return render(c, layout.Error(err))
	}
	return c.Redirect(302, "/admin/articles")
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func (h *AdminHandler) AdminNewArticle(c echo.Context) error {
	return render(c, admin.NewArticle())
}

func (h *AdminHandler) AdminEditArticleSave(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	var article admin.Article
	if err := c.Bind(&article); err != nil {
		return err
	}
	articleRow := repo.ArticleRow{
		ID:           id,
		Slug:         article.Slug,
		Title:        article.Title,
		Subtitle:     article.Subtitle,
		Body:         article.Content,
		MarkdownBody: article.Markdown,
		Image:        article.Image,
	}
	err = h.ArticleRepo.UpdateArticle(&articleRow)
	if err != nil {
		return err
	}
	return c.Redirect(302, "/admin/articles")
}

func (h *AdminHandler) AdminEditArticle(c echo.Context) error {
	idStr := c.Param("id")
	fmt.Println(c.FormParams())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	articleRow, err := h.ArticleRepo.GetArticleByID(id)
	if err != nil {
		return err
	}
	a := articleRow.ToAdminArticle()
	return render(c, admin.ArticleRowEdit(a))
}

func (h *AdminHandler) AdminLoginShow(c echo.Context) error {
	return render(c, admin.Login())
}

type AdminLogin struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (h *AdminHandler) AdminLogin(c echo.Context) error {
	var login AdminLogin
	if err := c.Bind(&login); err != nil {
		return c.Redirect(302, "/")
	}
	password, err := h.AdminRepo.Get(login.Username)
	if err != nil {
		return c.Redirect(302, "/admin")
	}
	if password != login.Password {
		return c.Redirect(302, "/admin")
	}
	cookie := &http.Cookie{
		Name:  "admin",
		Value: login.Username,
	}
	c.SetCookie(cookie)
	return c.Redirect(302, "/admin/articles")
}

func (h *AdminHandler) AdminArticlesShow(c echo.Context) error {
	articleRows, err := h.ArticleRepo.GetPublishedArticles()
	if err != nil {
		return err
	}
	articles := make([]admin.Article, len(articleRows))
	for pos, a := range articleRows {
		articles[pos] = a.ToAdminArticle()
	}
	return render(c, admin.Articles(articles))
}

func (h *AdminHandler) AdminUpdateArticle(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	var article admin.Article
	if err := c.Bind(&article); err != nil {
		return err
	}
	articleRow := repo.ArticleRow{
		ID:           id,
		Slug:         article.Slug,
		Title:        article.Title,
		Subtitle:     article.Subtitle,
		Body:         article.Content,
		MarkdownBody: article.Markdown,
		Image:        article.Image,
	}
	err = h.ArticleRepo.UpdateArticle(&articleRow)
	if err != nil {
		return err
	}
	return nil
}
