package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/graphql-go/graphql"
    "github.com/graphql-go/handler"
)

// Custom directive to handle substring logic
func substringDirective(value string, start int, length int) (string, error) {
    // Apply the substring operation
    if start < 0 || start > len(value) {
        return "", fmt.Errorf("start index out of bounds")
    }
    end := start + length
    if end > len(value) {
        end = len(value)
    }

    return value[start:end], nil
}

func main() {
    // Define the fields
    fields := graphql.Fields{
        "exampleField": &graphql.Field{
            Type: graphql.String,
            Args: graphql.FieldConfigArgument{
                "start": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
                "length": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                value := "Hello, GraphQL!"
                start := p.Args["start"].(int)
                length := p.Args["length"].(int)

                return substringDirective(value, start, length)
            },
        },
    }

    // Define the query type
    rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
    schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

    // Create the schema
    schema, err := graphql.NewSchema(schemaConfig)
    if err != nil {
        log.Fatalf("Failed to create new schema, error: %v", err)
    }

    // Create the GraphQL handler
    h := handler.New(&handler.Config{
        Schema:   &schema,
        Pretty:   true,
        GraphiQL: true,
    })

    // Serve the handler
    http.Handle("/graphql", h)
    fmt.Println("Now server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
