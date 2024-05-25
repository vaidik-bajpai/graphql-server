package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vaidik-bajpai/hackernews/database"
	"github.com/vaidik-bajpai/hackernews/graph"
	"github.com/vaidik-bajpai/hackernews/internal/data"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.OpenDB()
	defer database.CloseDB()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", data.Middleware(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", data.Middleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
