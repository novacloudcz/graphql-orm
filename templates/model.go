package templates

var Model = `package gen

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type {{.Name}}ResultType struct {
}

func (t *{{.Name}}ResultType) Items(ctx context.Context) (items []Todo) {
	rezCtx := graphql.GetResolverContext(ctx)
	//fields := graphql.CollectFieldsCtx(ctx, nil)
	//fmt.Println(rezCtx.Args, rezCtx.Parent.Args, fields)

	db := ctx.Value(DBContextKey).(*DB)
	q := db.Query()

	limit := rezCtx.Parent.Args["limit"]
	if value, ok := limit.(*int); ok && value != nil {
		q = q.Limit(*value)
	}
	offset := rezCtx.Parent.Args["offset"]
	if value, ok := offset.(*int); ok && value != nil {
		q = q.Offset(*value)
	}
	err := q.Model(&{{.Name}}{}).Find(&items).Error
	if err != nil {
		panic(err)
	}
	return
}

func (t *{{.Name}}ResultType) Count(ctx context.Context) (count int) {
	db := ctx.Value(DBContextKey).(*DB)
	err := db.Query().Model(&{{.Name}}{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	return
}`
