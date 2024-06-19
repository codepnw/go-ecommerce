package config

import (
	"fmt"
	"strconv"
)

type DbConfig interface {
	Url() string
	MaxOpenConns() int
	Driver() string
}

type db struct {
	driver         string
	host           string
	port           int
	protocal       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (c *config) Db() DbConfig {
	return c.db
}

func (d *db) Url() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.driver,
		d.username,
		d.password,
		d.host,
		strconv.Itoa(d.port),
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxOpenConns() int { return d.maxConnections }
func (d *db) Driver() string    { return d.driver }
