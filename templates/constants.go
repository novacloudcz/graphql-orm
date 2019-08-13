package templates

var Constants = `package gen

import (
)

type key int

const (
	KeyPrincipalID key = iota
	KeyJWTClaims key = iota
	RawSchema string = ` + "`{{.RawSchema}}`" + `
)
`
