package src

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
)

func New(db *gen.DB, ec *events.EventController) *Resolver {
	resolver := NewResolver(db, ec)

	resolver.Handlers.OnEvent = func(ctx context.Context, r *gen.GeneratedResolver, e *events.Event) (err error) {
		if e.Entity == "User" && e.Type == events.EventTypeUpdated {
			change := e.Change("firstName")
			if change != nil {
				var firstName string
				change.NewValueAs(&firstName)
				_, err = r.Handlers.CreateTask(ctx, r, map[string]interface{}{
					"title":      fmt.Sprintf("Hello %s!", firstName),
					"assigneeId": e.EntityID,
				})
			}
		}
		return nil
	}

	resolver.Handlers.CompanyReviews = func(ctx context.Context, r *gen.GeneratedCompanyResolver, obj *gen.Company) (res []*gen.Review, err error) {
		return []*gen.Review{
			&gen.Review{
				ID: "1",
			}, &gen.Review{
				ID: "2",
			},
		}, nil
	}

	resolver.Handlers.UserAddress = func(ctx context.Context, r *gen.GeneratedUserResolver, obj *gen.User) (res *gen.Address, err error) {
		if obj.AddressRaw != nil {
			res = &gen.Address{}
			err = json.Unmarshal([]byte(*obj.AddressRaw), res)
		}
		return
	}

	return resolver
}

func (r *QueryResolver) Hello(ctx context.Context) (string, error) {
	return "world", nil
}

func (r *QueryResolver) TopCompanies(ctx context.Context) (items []*gen.Company, err error) {
	err = r.DB.Query().Model(&gen.Company{}).Find(&items).Error
	return
}

func (r *CompanyResolver) UppercaseName(ctx context.Context, obj *gen.Company) (string, error) {
	name := ""
	if obj.Name != nil {
		name = *obj.Name
	}
	return strings.ToUpper(name), nil
}

func (r *ReviewResolver) Company(ctx context.Context, obj *gen.Review) (*gen.Company, error) {
	opts := gen.QueryCompanyHandlerOptions{
		ID: &obj.ReferenceID,
	}
	return r.Handlers.QueryCompany(ctx, r.GeneratedResolver, opts)
}

func (r *PlainEntityResolver) ShortText(ctx context.Context, obj *gen.PlainEntity) (string, error) {
	val := ""
	if obj.Text != nil {
		val = *obj.Text
	}
	return val, nil
}
