package templates

var HTTPHandler = `package gen

import (
	"context"
	"net/http"
	"strings"


	"github.com/99designs/gqlgen/handler"
	jwtgo "github.com/dgrijalva/jwt-go"
)

func GetHTTPServeMux(r ResolverRoot, db *DB) *http.ServeMux {
	mux := http.NewServeMux()

	gqlHandler := handler.GraphQL(NewExecutableSchema(Config{Resolvers: r}))

	loaders := GetLoaders(db)

	playgroundHandler := handler.Playground("GraphQL playground", "/api/graphql")
	mux.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		principalID := getPrincipalID(req)
		ctx := context.WithValue(req.Context(), KeyPrincipalID, principalID)
		ctx = context.WithValue(ctx, "loaders", loaders)
		req = req.WithContext(ctx)
		if req.Method == "GET" {
			playgroundHandler(res, req)
		} else {
			gqlHandler(res, req)
		}
	})
	handler := mux

	return handler
}

func getPrincipalIDFromContext(ctx context.Context) *string {
	v, _ := ctx.Value(KeyPrincipalID).(*string)
	return v
}
func getJWTClaimsFromContext(ctx context.Context) *JWTClaims {
	v, _ := ctx.Value(KeyJWTClaims).(*JWTClaims)
	return v
}

func getPrincipalID(req *http.Request) *string {
	pID := req.Header.Get("principal-id")
	if pID != "" {
		return &pID
	}
	c, _ := getJWTClaims(req)
	if c == nil {
		return nil
	}
	return &c.Subject
}

type JWTClaims struct {
	jwtgo.StandardClaims
	Scope *string
}

func getJWTClaims(req *http.Request) (*JWTClaims, error) {
	var p *JWTClaims

	tokenStr := strings.Replace(req.Header.Get("authorization"), "Bearer ", "", 1)
	if tokenStr == "" {
		return p, nil
	}

	p = &JWTClaims{}
	jwtgo.ParseWithClaims(tokenStr, p, nil)
	return p, nil
}

`
