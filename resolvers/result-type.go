package resolvers

import (
	"context"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/jinzhu/gorm"
)

type EntityFilter interface {
	Apply(db *gorm.DB) (*gorm.DB, error)
}
type EntitySort interface {
	String() string
}

type EntityResultType struct {
	Offset *int
	Limit  *int
	Query  *string
	Sort   []EntitySort
	Filter EntityFilter
}

// GetResultTypeItems ...
func (r *EntityResultType) GetItems(ctx context.Context, db *gorm.DB, out interface{}) error {
	q := db

	if r.Limit != nil {
		q = q.Limit(*r.Limit)
	}
	if r.Offset != nil {
		q = q.Offset(*r.Offset)
	}

	for _, s := range r.Sort {
		direction := "ASC"
		_s := s.String()
		if strings.HasSuffix(_s, "_DESC") {
			direction = "DESC"
		}
		col := strcase.ToLowerCamel(strings.ToLower(strings.TrimSuffix(_s, "_"+direction)))
		q = q.Order(col + " " + direction)
	}

	q, err := r.Filter.Apply(q)
	if err != nil {
		return err
	}

	return q.Find(out).Error
}

// GetCount ...
func (r *EntityResultType) GetCount(ctx context.Context, db *gorm.DB, out interface{}) (count int, err error) {
	err = db.Model(out).Count(&count).Error
	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
