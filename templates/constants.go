package templates

var Constants = `package gen

import (
)

type key int

const (
	KeyPrincipalID      key    = iota
	KeyLoaders          key    = iota
	KeyExecutableSchema key    = iota
	KeyJWTClaims        key    = iota
	SchemaSDL string = ` + "`{{.SchemaSDL}}`" + `
)
`
