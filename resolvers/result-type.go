package resolvers

import (
	"context"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/jinzhu/gorm"
)

type EntityFilter interface {
	Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntitySort interface {
	String() string
}

type EntityResultType struct {
	Offset *int
	Limit  *int
	Query  EntityFilterQuery
	Sort   []EntitySort
	Filter EntityFilter
}

// GetResultTypeItems ...
func (r *EntityResultType) GetItems(ctx context.Context, db *gorm.DB, alias string, out interface{}) error {
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

	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err := r.Query.Apply(ctx, &wheres, &values, &joins)
	if err != nil {
		return err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, &wheres, &values, &joins)
		if err != nil {
			return err
		}
	}

	if len(wheres) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), values...)
	}

	uniqueJoins := map[string]bool{}
	for _, join := range joins {
		uniqueJoins[join] = true
	}

	for join := range uniqueJoins {
		q = q.Joins(join)
	}

	q = q.Group(alias + ".id")
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
