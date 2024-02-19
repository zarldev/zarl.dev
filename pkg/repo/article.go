package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zarldev/zarldotdev/view/admin"
	"github.com/zarldev/zarldotdev/view/article"
)

type ArticleRepository struct {
	conn *Connection
}

func NewArticleRepository(conn *Connection) (*ArticleRepository, error) {
	ar := &ArticleRepository{conn}
	err := ar.createTable()
	if err != nil {
		return nil, err
	}
	return ar, nil
}

type ArticleRow struct {
	ID           int
	Slug         string
	Title        string
	Subtitle     string
	Body         string
	MarkdownBody string
	Created      string
	Updated      string
	Image        string
	Published    bool
	Claps        sql.NullInt32
}

func (a ArticleRow) ToArticle() article.Article {
	return article.Article{
		ID:       a.ID,
		Title:    a.Title,
		Subtitle: a.Subtitle,
		Content:  a.Body,
		Image:    a.Image,
		Slug:     a.Slug,
		Claps:    claps(a),
		Created:  timeFrom(a.Created),
	}
}

func (a ArticleRow) ToAdminArticle() admin.Article {
	return admin.Article{
		ID:       a.ID,
		Title:    a.Title,
		Subtitle: a.Subtitle,
		Content:  a.Body,
		Image:    a.Image,
		Slug:     a.Slug,
		Created:  timeFrom(a.Created),
		Updated:  timeFrom(a.Updated),
		Markdown: a.MarkdownBody,
		Claps:    claps(a),
	}
}

func claps(a ArticleRow) int {
	claps := 0
	if a.Claps.Valid {
		claps = int(a.Claps.Int32)
	}
	return claps
}

func timeFrom(timeStr string) time.Time {
	created, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("failed to parse time", err)
		created = time.Time{}
	}
	return created
}

func (r *ArticleRepository) createTable() error {
	_, err := r.conn.write.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id serial PRIMARY KEY NOT NULL,
			slug text NOT NULL,
			title text NOT NULL,
			subtitle text NOT NULL,
			body text NOT NULL,
			markdown_body text NOT NULL,
			image text,
			created timestamp default current_timestamp NOT NULL,
			updated timestamp default current_timestamp NOT NULL,
			published boolean default false NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (r *ArticleRepository) CreateArticle(a *ArticleRow) error {
	_, err := r.conn.write.Exec(`
		INSERT INTO articles (slug, title, subtitle, body, markdown_body, image, published)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, a.Slug, a.Title, a.Subtitle, a.Body, a.MarkdownBody, a.Image, a.Published)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

func (r *ArticleRepository) GetArticleBySlug(slug string) (*ArticleRow, error) {
	row := r.conn.read.QueryRow(`SELECT a.id, a.slug, a.title, a.subtitle, a.body, a.markdown_body, a.created, a.updated, a.image, a.published, c.count AS claps
	FROM articles AS a
	LEFT JOIN claps AS c ON a.id = c.article_id
	WHERE a.slug = $1 ORDER BY a.id DESC`, slug)
	return rowToArticleRow(row)
}

func rowToArticleRow(row *sql.Row) (*ArticleRow, error) {
	a := &ArticleRow{}
	err := row.Scan(&a.ID, &a.Slug, &a.Title, &a.Subtitle, &a.Body, &a.MarkdownBody, &a.Created, &a.Updated, &a.Image, &a.Published, &a.Claps)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return a, nil
}

func rowsToArticleRow(rows *sql.Rows) (*ArticleRow, error) {
	a := &ArticleRow{}
	err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Subtitle, &a.Body, &a.MarkdownBody, &a.Created, &a.Updated, &a.Image, &a.Published, &a.Claps)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return a, nil
}

func (r *ArticleRepository) GetArticleByID(id int) (*ArticleRow, error) {
	row := r.conn.read.QueryRow(`SELECT a.id, a.slug, a.title, a.subtitle, a.body, a.markdown_body, a.created, a.image, a.published c.count AS claps
	FROM articles AS a
	LEFT JOIN claps AS c ON a.id = c.article_id
	WHERE a.id = $1 ORDER BY a.id DESC`, id)
	return rowToArticleRow(row)
}

func (r *ArticleRepository) GetPublishedArticles() ([]*ArticleRow, error) {
	rows, err := r.conn.read.Query(`
	SELECT a.id, a.slug, a.title, a.subtitle, a.body, a.markdown_body, a.created, a.updated, a.image, a.published, c.count AS claps
	FROM articles AS a 
	LEFT JOIN claps AS c ON a.id = c.article_id WHERE a.published = 'true'
	ORDER BY a.id DESC;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	articles, err := articleRowsFromRows(rows)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func articleRowsFromRows(rows *sql.Rows) ([]*ArticleRow, error) {
	var articles []*ArticleRow
	for rows.Next() {
		a, err := rowsToArticleRow(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil
}

func (r *ArticleRepository) UpdateArticle(a *ArticleRow) error {
	_, err := r.conn.write.Exec(`
		UPDATE articles SET slug = $1, title = $2, subtitle = $3, body = $4, markdown_body = $5, image = $6, updated = current_timestamp 
		WHERE id = $7;
	`, a.Slug, a.Title, a.Subtitle, a.Body, a.MarkdownBody, a.Image, a.ID)
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}
	return nil
}
