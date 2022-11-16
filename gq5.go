package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

type employee struct {
	ID   int    `json:"employeeId"`
	Name string `json:"employeeName"`
}

var e = []employee{{1, "John"}, {2, "Jay"}, {3, "Smith"}, {4, "Ricky"}}

var schema graphql.Schema

var employeeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "employee",
		Fields: graphql.Fields{
			"employeeId": &graphql.Field{
				Type: graphql.Int,
			},
			"employeeName": &graphql.Field{
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

var rootMutationFields = graphql.Fields{
	"createEmployee": &graphql.Field{
		Type:        employeeType,
		Description: "Create a new employees",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok1 := p.Args["id"].(int)
			name, ok2 := p.Args["name"].(string)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("argument type mismatch")
			}
			for _, v := range e {
				if v.ID == id {
					return nil, fmt.Errorf("id already present")
				}
			}
			e = append(e, employee{ID: id, Name: name})
			return employee{ID: id, Name: name}, nil
		},
	},
}

func main() {
	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: rootQueryFields,
	}
	rootMutation := graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: rootMutationFields,
	}

	schemaConfig := graphql.SchemaConfig{
		Mutation: graphql.NewObject(rootMutation),
		Query:    graphql.NewObject(rootQuery),
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
