package src

import (
	"context"

	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
)

func New(db *gen.DB, ec *events.EventController) *gen.GeneratedResolver {
	resolver := gen.NewResolver(db, ec)

	// resolver.Handlers.CompanyReview = func(ctx context.Context, r *gen.GeneratedCompanyResolver, obj *gen.Company) (res *gen.Review, err error) {
	// 	id := "dummy_ID"
	// 	return &gen.Review{
	// 		ID: id,
	// 	}, nil
	// }

	resolver.Handlers.CompanyReviews = func(ctx context.Context, r *gen.GeneratedCompanyResolver, obj *gen.Company) (res []*gen.Review, err error) {
		return []*gen.Review{
			&gen.Review{
				ID: "1",
			}, &gen.Review{
				ID: "2",
			},
		}, nil
	}

	resolver.Handlers.ReviewCompany = func(ctx context.Context, r *gen.GeneratedReviewResolver, obj *gen.Review) (res *gen.Company, err error) {
		// fmt.Println("??", obj.ID, obj.ReferenceID)
		return r.Query().Company(ctx, &obj.ReferenceID, nil, nil)
	}

	return resolver
}
