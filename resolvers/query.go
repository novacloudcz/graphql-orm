package resolvers

import (
	"context"

	"github.com/jinzhu/gorm"
)

// GetItem ...
func GetItem(ctx context.Context, db *gorm.DB, out interface{}, id *string) error {
	return db.Find(out, "id = ?", id).Error
}
