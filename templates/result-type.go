package templates

var ResultType = `package gen

import (
	"context"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"

	"github.com/jinzhu/gorm"
)

func GetItem(ctx context.Context, db *gorm.DB, out interface{}, id *string) error {
	return db.Find(out, "id = ?", id).Error
}

func GetItemForRelation(ctx context.Context, db *gorm.DB, obj interface{}, relation string, out interface{}) error {
	return db.Model(obj).Related(out, relation).Error
}

type EntityFilter interface {
	Apply(ctx context.Context, dialect gorm.Dialect, wheres *[]string, whereValues *[]interface{}, havings *[]string, havingValues *[]interface{}, joins *[]string) error
}
type EntityFilterQuery interface {
	Apply(ctx context.Context, dialect gorm.Dialect, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error
}


type SortInfo struct {
	Field         string
	Direction     string
	IsAggregation bool
}
func (si *SortInfo) String() string {
	return fmt.Sprintf("%s %s", si.Field, si.Direction)
}
type EntitySort interface {
	Apply(ctx context.Context, dialect gorm.Dialect, sorts *[]SortInfo, joins *[]string) error
}

type EntityResultType struct {
	Offset       *int
	Limit        *int
	Query        EntityFilterQuery
	Sort         []EntitySort
	Filter       EntityFilter
	Fields       []*ast.Field
	SelectionSet *ast.SelectionSet
}

type GetItemsOptions struct {
	Alias      string
	Preloaders []string
}

// GetResultTypeItems ...
func (r *EntityResultType) GetItems(ctx context.Context, db *gorm.DB, opts GetItemsOptions, out interface{}) error {
	subq := db.Model(out)
	q := db
	subqGroups := []string{opts.Alias + ".id"}
	subqFields := []string{"DISTINCT(" + opts.Alias + ".id) as id"}
	qSorts := []string{}
	subqSorts := []string{}

	if r.Limit != nil {
		// q = q.Limit(*r.Limit)
		subq = subq.Limit(*r.Limit)
	}
	if r.Offset != nil {
		// q = q.Offset(*r.Offset)
		subq = subq.Offset(*r.Offset)
	}

	dialect := q.Dialect()

	wheres := []string{}
	havings := []string{}
	whereValues := []interface{}{}
	havingValues := []interface{}{}
	joins := []string{}
	sorts := []SortInfo{}

	err := r.Query.Apply(ctx, dialect, r.SelectionSet, &wheres, &whereValues, &joins)
	if err != nil {
		return err
	}

	for _, sort := range r.Sort {
		sort.Apply(ctx, dialect, &sorts, &joins)
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &whereValues, &havings, &havingValues, &joins)
		if err != nil {
			return err
		}
	}

	if len(sorts) > 0 {
		for i, sort := range sorts {
			if !sort.IsAggregation {
				subqGroups = append(subqGroups, fmt.Sprintf("%s", sort.Field))
			}
			subqFields = append(subqFields, fmt.Sprintf("%s as "+dialect.Quote("sort_key_%d"), sort.Field, i))
			qSorts = append(qSorts, fmt.Sprintf(dialect.Quote("filter_table")+"."+dialect.Quote("sort_key_%d")+" %s", i, sort.Direction))
			subqSorts = append(subqSorts, fmt.Sprintf(dialect.Quote("sort_key_%d")+" %s", i, sort.Direction))
		}
	}
	if len(wheres) > 0 {
		subq = subq.Where(strings.Join(wheres, " AND "), whereValues...)
	}
	if len(havings) > 0 {
		subq = subq.Having(strings.Join(havings, " AND "), havingValues...)
	}

	uniqueJoinsMap := map[string]bool{}
	uniqueJoins := []string{}
	for _, join := range joins {
		if !uniqueJoinsMap[join] {
			uniqueJoinsMap[join] = true
			uniqueJoins = append(uniqueJoins, join)
		}
	}

	for _, join := range uniqueJoins {
		q = q.Joins(join)
		subq = subq.Joins(join)
	}

	if len(opts.Preloaders) > 0 {
		for _, p := range opts.Preloaders {
			q = q.Preload(p)
		}
	}
	subq = subq.Group(strings.Join(subqGroups, ", ")).Select(strings.Join(subqFields, ", "))
	subq = subq.Order(strings.Join(subqSorts, ", "))
	q = q.Order(strings.Join(qSorts, ", "))
	q = q.Joins("INNER JOIN (?) as filter_table ON filter_table.id = "+opts.Alias+".id", subq.QueryExpr())

	return q.Find(out).Error
}

// GetCount ...
func (r *EntityResultType) GetCount(ctx context.Context, db *gorm.DB, opts GetItemsOptions, out interface{}) (count int, err error) {
	q := db

	dialect := q.Dialect()
	wheres := []string{}
	havings := []string{}
	whereValues := []interface{}{}
	havingValues := []interface{}{}
	joins := []string{}

	err = r.Query.Apply(ctx, dialect, r.SelectionSet, &wheres, &whereValues, &joins)
	if err != nil {
		return 0, err
	}

	if r.Filter != nil {
		err = r.Filter.Apply(ctx, dialect, &wheres, &whereValues, &havings, &havingValues, &joins)
		if err != nil {
			return 0, err
		}
	}

	if len(wheres) > 0 {
		q = q.Where(strings.Join(wheres, " AND "), whereValues...)
	}
	if len(havings) > 0 {
		q = q.Having(strings.Join(havings, " AND "), havingValues...)
	}

	uniqueJoinsMap := map[string]bool{}
	uniqueJoins := []string{}
	for _, join := range joins {
		if !uniqueJoinsMap[join] {
			uniqueJoinsMap[join] = true
			uniqueJoins = append(uniqueJoins, join)
		}
	}

	for _, join := range uniqueJoins {
		q = q.Joins(join)
	}
	q = q.Model(out).Group(opts.Alias + ".id")
	err = db.Model(out).Joins("INNER JOIN (?) as filter_table ON filter_table.id = "+opts.Alias+".id", q.QueryExpr()).Count(&count).Error

	return
}

func (r *EntityResultType) GetSortStrings() []string {
	return []string{}
}
`
