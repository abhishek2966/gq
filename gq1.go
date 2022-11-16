package main1

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

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

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Panic(err)
	}

	requestQuery := `
		{
			msg
		}
	`

	params := graphql.Params{
		Schema:        schema,
		RequestString: requestQuery,
	}

	res := graphql.Do(params)
	if len(res.Errors) != 0 {
		fmt.Println(res.Errors)
		return
	}

	resJson, _ := json.Marshal(res)
	fmt.Println(string(resJson))
}
