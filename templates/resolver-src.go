package templates

var ResolverSrc = `package src

import (
	"{{.Config.Package}}/gen"
	"github.com/novacloudcz/graphql-orm/events"
)


func New(db *gen.DB, ec *events.EventController) *gen.GeneratedResolver {
	resolver := gen.NewResolver(db, ec)

	// resolver.Handlers.CreateUser = func(ctx context.Context, r *gen.GeneratedMutationResolver, input map[string]interface{}) (item *gen.Company, err error) {
	// 	return gen.CreateUserHandler(ctx, r, input)
	// }

	return resolver
}
`
