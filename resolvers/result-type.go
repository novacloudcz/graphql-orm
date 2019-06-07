package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jinzhu/gorm"
)

// GetResultTypeItems
func GetResultTypeItems(ctx context.Context, db *gorm.DB, out interface{}) error {
	rezCtx := graphql.GetResolverContext(ctx)
	//fields := graphql.CollectFieldsCtx(ctx, nil)
	//fmt.Println(rezCtx.Args, rezCtx.Parent.Args, fields)

	// db := ctx.Value(DBContextKey).(*DB)
	q := db

	limit := rezCtx.Parent.Args["limit"]
	if value, ok := limit.(*int); ok && value != nil {
		q = q.Limit(*value)
	}
	offset := rezCtx.Parent.Args["offset"]
	if value, ok := offset.(*int); ok && value != nil {
		q = q.Offset(*value)
	}
	return q.Find(out).Error
}

// GetResultTypeCount ...
func GetResultTypeCount(ctx context.Context, db *gorm.DB, out interface{}) (count int, err error) {
	err = db.Model(out).Count(&count).Error
	return
}
