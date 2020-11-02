package templates

var Database = `package gen

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB ...
type DB struct {
	db *gorm.DB
}

// NewDBFromEnvVars Create database client using DATABASE_URL environment variable
func NewDBFromEnvVars() *DB {
	urlString := os.Getenv("DATABASE_URL")
	if urlString == "" {
		panic(fmt.Errorf("missing DATABASE_URL environment variable"))
	}
	return NewDBWithString(urlString)
}

func TableName(name string) string {
	prefix := os.Getenv("TABLE_NAME_PREFIX")
	if prefix != "" {
		return prefix + "_" + name
	}
	return name
}

// NewDB ...
func NewDB(db *gorm.DB) *DB {

	v := DB{db}
	return &v
}

// NewDBWithString creates database instance with database URL string
func NewDBWithString(urlString string) *DB {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	urlString = getConnectionString(u)

	var dialector gorm.Dialector
	switch u.Scheme {
	case "sqlite3":
		dialector = sqlite.Open(urlString)
	case "mysql":
		dialector = mysql.Open(urlString)
	case "postgres":
		dialector = postgres.Open(urlString)
	}

	prefix := os.Getenv("TABLE_NAME_PREFIX")
	if prefix != "" {
		prefix += "_"
	}

	logMode := logger.Silent
	if os.Getenv("DEBUG") == "true" {
		logMode = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: prefix,
		},
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	}

	if urlString == "sqlite3://:memory:" {
		rawDB, _ := db.DB()
		rawDB.SetMaxIdleConns(1)
		rawDB.SetConnMaxLifetime(time.Second * 300)
		rawDB.SetMaxOpenConns(1)
	} else {
		rawDB, _ := db.DB()
		rawDB.SetMaxIdleConns(5)
		rawDB.SetConnMaxLifetime(time.Second * 60)
		rawDB.SetMaxOpenConns(10)
	}

	return NewDB(db)
}

func getConnectionString(u *url.URL) string {
	if u.Scheme == "postgres" {
		password, _ := u.User.Password()
		params := u.Query()
		params.Set("host", strings.Split(u.Host, ":")[0])
		params.Set("port", u.Port())
		params.Set("user", u.User.Username())
		params.Set("password", password)
		params.Set("dbname", strings.TrimPrefix(u.Path, "/"))
		return strings.Replace(params.Encode(), "&", " ", -1)
		// return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, u.Port(), u.User.Username(), password, strings.TrimPrefix(u.Path, "/"))
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

// AutoMigrate run basic gorm automigration
func (db *DB) AutoMigrate() error {
	return AutoMigrate(db.db)
}

// Migrate run migrations using automigrate
func (db *DB) Migrate(migrations []*gormigrate.Migration) error {
	options := gormigrate.DefaultOptions
	options.TableName = TableName("migrations")
	return Migrate(db.db, options, migrations)
}

// Close ...
func (db *DB) Close() error {
	rawDB, _ := db.db.DB()
	return rawDB.Close()
}

// Ping ...
func (db *DB) Ping() error {
	rawDB, _ := db.db.DB()
	return rawDB.Ping()
}
`
