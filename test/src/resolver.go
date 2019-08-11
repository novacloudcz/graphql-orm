package src

import (
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
)

type Resolver struct {
	*gen.GeneratedResolver
}

func New(db *gen.DB, ec *events.EventController) *Resolver {
	return &Resolver{&gen.GeneratedResolver{db, ec}}
}

// This is example how to override default resolver to provide customizations

// 1) Create resolver for specific part of the query (mutation, query, result types etc.)
// type MutationResolver struct{ *gen.GeneratedMutationResolver }

// 2) Override Resolver method for returning your own resolver
// func (r *Resolver) Mutation() gen.MutationResolver {
// 	return &MutationResolver{&gen.GeneratedMutationResolver{r.GeneratedResolver}}
// }

// 3) Implement custom logic for your resolver
// Replace XXX with your entity name (you can find definition of these methods in generated resolvers)
// func (r *MutationResolver) CreateXXX(ctx context.Context, input map[string]interface{}) (item *gen.Company, err error) {
//	// example call of your own logic
//  if err := validateCreateXXXInput(input); err != nil {
// 		return nil, err
//	}
// 	return r.GeneratedMutationResolver.CreateXXX(ctx, input)
// }
