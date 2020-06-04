package resolvers

import "github.com/99designs/gqlgen/graphql"

func getServer(rctx *graphql.FieldContext) (string, bool) {
	server := ""
	ok := false
	parent := rctx.Parent
	for parent != nil {
		server, ok = parent.Args["server"].(string)
		if ok {
			break
		}
		parent = parent.Parent
	}
	return server, ok
}
