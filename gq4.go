package main1

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
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

	h := handler.New(&handler.Config{
		Schema: &schema,
		// Pretty print JSON response
		Pretty: true,
		// Host a GraphiQL Playground to use for testing Queries
		GraphiQL:   true,
		Playground: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
