package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/LFSCamargo/twitter-go/auth"
	"github.com/LFSCamargo/twitter-go/database"
	"github.com/LFSCamargo/twitter-go/graph"
	"github.com/LFSCamargo/twitter-go/graph/generated"
	"github.com/go-chi/chi"
	"github.com/go-chi/httplog"
)

const defaultPort = "8000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logger := httplog.NewLogger("twitter-go-logging", httplog.Options{
		JSON:    true,
		Concise: true,
	})

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	database.Connect()

	router := chi.NewRouter()
	router.Use(httplog.RequestLogger(logger))
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
