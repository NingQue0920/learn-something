//go:build clickhouse
// +build clickhouse

package middleware

import (
	"database/sql"
)

type ClickHouse struct {
	conn *sql.DB
}

func init() {
	RegisterMiddleware("clickhouse", NewClickHouse)
}

func NewClickHouse() (Middleware, error) {
	return &ClickHouse{}, nil
}

func (c *ClickHouse) Initialize() error {
	// init clickhouse
	var err error
	c.conn, err = sql.Open("clickhouse", "tcp://localhost:9000?debug=true")
	_, err = clickhouse.Open("tcp://localhost:9000?debug=true")
	return err
}
func (c *ClickHouse) Read(query string) (any, error) {
	return c.conn.Query(query)

}
