package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/LFSCamargo/graphql-go-boilerplate/auth"
	"github.com/LFSCamargo/graphql-go-boilerplate/database"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph"
	"github.com/LFSCamargo/graphql-go-boilerplate/graph/generated"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	database.Connect()

	router := chi.NewRouter()

	router.Use(auth.Middleware())
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("Server exposed at http://localhost:%s/graphql", port)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		panic(err)
	}
}
