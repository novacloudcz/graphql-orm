package resolvers

import (
	"context"

	"github.com/mitchellh/mapstructure"

	"github.com/jinzhu/gorm"
)

func CreateItem(ctx context.Context, db *gorm.DB, model interface{}, data map[string]interface{}) error {
	if err := mapstructure.Decode(data, model); err != nil {
		return err
	}

	return db.Create(model).Error
}

func UpdateItem(ctx context.Context, db *gorm.DB, model interface{}, data map[string]interface{}) error {
	if err := mapstructure.Decode(data, model); err != nil {
		return err
	}
	return db.Save(model).Error
}

func DeleteItem(ctx context.Context, db *gorm.DB, model interface{}, id string) error {
	return db.Delete(model, "id = ?", id).Error
}
