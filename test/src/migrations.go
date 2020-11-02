package src

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/novacloudcz/graphql-orm/test/gen"
	"gorm.io/gorm"
)

func GetMigrations(db *gen.DB) []*gormigrate.Migration {
	return []*gormigrate.Migration{
		&gormigrate.Migration{
			ID: "INIT",
			Migrate: func(tx *gorm.DB) error {
				return db.AutoMigrate()
			},
			Rollback: func(tx *gorm.DB) error {
				// there's not much we can do if initialization/automigration failes
				return nil
			},
		},
		&gormigrate.Migration{
			ID: "01create_user_blah",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					Blah string
				}
				return tx.AutoMigrate(&User{})
			},
		},
		&gormigrate.Migration{
			ID: "02drop_user_blah",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					blah string
				}
				return tx.Model(User{}).DropColumn("blah").Error
			},
		},
	}
}
