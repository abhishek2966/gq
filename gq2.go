package main1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

func main() {
	fields := graphql.Fields{
		"msg": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "Welcome", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	schema, _ = graphql.NewSchema(schemaConfig)

	r := mux.NewRouter()
	r.HandleFunc("/graphql", executeQuery)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
func executeQuery(w http.ResponseWriter, r *http.Request) {
	requestQuery := r.URL.Query().Get("query")
	params := graphql.Params{
		Schema:        schema,
		RequestString: requestQuery,
	}

	res := graphql.Do(params)
	if len(res.Errors) != 0 {
		http.Error(w, fmt.Sprint(res.Errors), http.StatusBadRequest)
		return
	}

	resJson, _ := json.Marshal(res)
	fmt.Fprintln(w, string(resJson))
}
