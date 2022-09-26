package main

import (
	"os"
	"sample/api/core"
	"sample/api/gorillamux"
	"sample/api/graph"
	"sample/api/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

var router *mux.Router
var muxAdapter *gorillamux.GorillaMuxAdapter

func init() {
	router = mux.NewRouter()
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})

	srv := handler.NewDefaultServer(schema)

	router.Handle("/query", srv)

	playgroundQueryPath := "/api/v1/web" + os.Getenv("__OW_ACTION_NAME") + "/query"
	router.Handle("/", playground.Handler("GraphQL Playground", playgroundQueryPath))

	muxAdapter = gorillamux.New(router)
}

func Main(args map[string]interface{}) map[string]interface{} {

	mainArgs := core.MainArgsFromMap(&args)
	resp, _ := muxAdapter.MainFnAdapter(mainArgs)

	data := core.MainArgsToMap(&resp)

	if data == nil {
		m := core.ErrorResponse(500)
		return core.MainArgsToMap(&m)
	}

	return data
}
