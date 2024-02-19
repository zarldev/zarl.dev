package repo

type ClapsRepository struct {
	conn *Connection
}

type Clap struct {
	ID    int
	Count int
}

func NewClapsRepository(conn *Connection) (*ClapsRepository, error) {
	cr := &ClapsRepository{conn: conn}
	err := cr.createTable()
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (c *ClapsRepository) createTable() error {
	_, err := c.conn.write.Exec("CREATE TABLE IF NOT EXISTS claps (article_id INTEGER PRIMARY KEY, count INTEGER)")
	if err != nil {
		return err
	}
	return nil
}

func (c *ClapsRepository) Get(id int) (int, error) {
	var count int
	err := c.conn.read.QueryRow("SELECT count FROM claps WHERE article_id = ?", id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *ClapsRepository) Increment(id int) (int, error) {
	_, err := c.conn.write.Exec("INSERT INTO claps (article_id, count) VALUES (?, 1) ON CONFLICT(article_id) DO UPDATE SET count = count + 1", id)
	if err != nil {
		return 0, err
	}
	count, err := c.Get(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
