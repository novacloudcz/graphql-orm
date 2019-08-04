package templates

var Database = `package gen

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB ...
type DB struct {
	db *gorm.DB
}

// NewDB ...
func NewDB(db *gorm.DB) *DB {
	prefix := os.Getenv("TABLE_NAME_PREFIX")
	if prefix != "" {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return prefix + "_" + defaultTableName
		}
	}
	v := DB{db}
	return &v
}

// NewDBWithString ...
func NewDBWithString(urlString string) *DB {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	urlString = getConnectionString(u)

	db, err := gorm.Open(u.Scheme, urlString)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns({{.Config.MaxIdleConnections}})
	db.DB().SetConnMaxLifetime(time.Seconds*{{.Config.ConnMaxLifetime}})
	db.DB().SetMaxOpenConns({{.Config.MaxOpenConnections}})
	db.LogMode(true)
	return NewDB(db)
}

func getConnectionString(u *url.URL) string {
	if u.Scheme == "postgres" {
		password, _ := u.User.Password()
		host := strings.Split(u.Host, ":")[0]
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, u.Port(), u.User.Username(), password, strings.TrimPrefix(u.Path, "/"))
	}
	if u.Scheme != "sqlite3" {
		u.Host = "tcp(" + u.Host + ")"
	}
	if u.Scheme == "mysql" {
		q := u.Query()
		q.Set("parseTime", "true")
		u.RawQuery = q.Encode()
	}
	return strings.Replace(u.String(), u.Scheme+"://", "", 1)
}

// Query ...
func (db *DB) Query() *gorm.DB {
	return db.db
}

// AutoMigrate ...
func (db *DB) AutoMigrate() {
	db.db.AutoMigrate({{range .Model.Objects}}
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
