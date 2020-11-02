package templates

var MigrationsSrc = `package src

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"
	"github.com/go-gormigrate/gormigrate/v2"
)

func GetMigrations(db *gen.DB) []*gormigrate.Migration {
	return []*gormigrate.Migration{
		&gormigrate.Migration{
			ID: "0000_INIT",
			Migrate: func(tx *gorm.DB) error {
				return db.AutoMigrate()
			},
			Rollback: func(tx *gorm.DB) error {
				// there's not much we can do if initialization/automigration failes
				return nil
			},
		},
	}
}
`
