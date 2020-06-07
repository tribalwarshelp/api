package resolvers

import "github.com/99designs/gqlgen/graphql"

func getServer(rctx *graphql.FieldContext) (string, bool) {
	server := ""
	ok := false
	for rctx != nil {
		server, ok = rctx.Args["server"].(string)
		if ok {
			break
		}
		rctx = rctx.Parent
	}
	return server, ok
}
