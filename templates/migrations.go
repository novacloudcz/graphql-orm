package templates

var Migrations = `package gen

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"
	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrate(db *gorm.DB, options *gormigrate.Options, migrations []*gormigrate.Migration) error {
	m := gormigrate.New(db, options, migrations)

	// // it's possible to use this, but in case of any specific keys or columns are created in migrations, they will not be generated by automigrate
	// m.InitSchema(func(tx *gorm.DB) error {
	// 	return AutoMigrate(db)
	// })

	return m.Migrate();
}

func AutoMigrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate({{range $obj := .Model.ObjectEntities}}
		{{.Name}}{},{{end}}
	)

	return
}
`
