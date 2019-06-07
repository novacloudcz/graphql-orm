package templates

var Database = `package gen

import (
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type key int

const (
	DBContextKey key = iota
)

// DB ...
type DB struct {
	db *gorm.DB
}

// NewDB ...
func NewDB(db *gorm.DB) *DB {
	v := DB{db}
	return &v
}

// NewDBWithString ...
func NewDBWithString(urlString string) *DB {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	if u.Scheme != "sqlite3" {
		u.Host = "tcp(" + u.Host + ")"
	}

	urlString = strings.Replace(u.String(), u.Scheme+"://", "", 1)

	db, err := gorm.Open(u.Scheme, urlString)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return NewDB(db)
}

// Query ...
func (db *DB) Query() *gorm.DB {
	return db.db
}

// AutoMigrate ...
func (db *DB) AutoMigrate() {
	db.db.AutoMigrate({{range .Objects}}
		{{.Name}}{},{{end}}
	)
}

// Close ...
func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Ping() error {
	return db.db.DB().Ping()
}
`
