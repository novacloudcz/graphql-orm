package src

import (
	"context"

	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
)

func New(db *gen.DB, ec *events.EventController) *gen.GeneratedResolver {
	resolver := gen.NewResolver(db, ec)

	resolver.Handlers.CompanyReview = func(ctx context.Context, r *gen.GeneratedCompanyResolver, obj *gen.Company) (res *gen.Review, err error) {
		id := "dummy_ID"
		return &gen.Review{
			ID: id,
		}, nil
	}

	return resolver
}
