package main1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

type employee struct {
	ID   int
	Name string
}

var e = []employee{{1, "John"}, {2, "Jay"}, {3, "Smith"}, {4, "Ricky"}}

var schema graphql.Schema

var employeeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "employee",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
var rootQueryFields = graphql.Fields{
	"employees": &graphql.Field{
		Type:        graphql.NewList(employeeType),
		Description: "Get all employees",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return e, nil
		},
	},
}

func main() {
	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: rootQueryFields,
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	schema, _ = graphql.NewSchema(schemaConfig)

	r := mux.NewRouter()
	r.HandleFunc("/graphql", executeQuery).Methods("GET")
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
