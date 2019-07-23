package resolvers

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/jinzhu/gorm"
)

type EntityFilter interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, values *[]interface{}, joins *[]string) error
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

	dialect := q.Dialect()
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err := r.Query.Apply(ctx, dialect, &wheres, &values, &joins)
	if err != nil {
		return err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &values, &joins)
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
	q := db

	dialect := q.Dialect()
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}

	err = r.Query.Apply(ctx, dialect, &wheres, &values, &joins)
	if err != nil {
		return 0, err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &values, &joins)
		if err != nil {
			return 0, err
		}
	}

	if len(wheres) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), values...)
	}

	uniqueJoinsMap := map[string]bool{}
	uniqueJoins := []string{}
	for _, join := range joins {
		if uniqueJoinsMap[join] == false {
			uniqueJoinsMap[join] = true
			uniqueJoins = append(uniqueJoins, join)
		}
	}

	fmt.Println("joins!!!", uniqueJoins)
	for _, join := range uniqueJoins {
		q = q.Joins(join)
	}
	err = q.Model(out).Count(&count).Error
	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
