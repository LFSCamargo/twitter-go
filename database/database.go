package database

import (
	"log"
	"os"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect - Is the function that connects to the mongoDB
func Connect() {
	mongouri := os.Getenv("MONGO_URI")
	mongodatabase := os.Getenv("MONGO_DATABASE")
	if mongouri == "" && mongodatabase == "" {
		mongouri = "mongodb://localhost:27017"
		mongodatabase = "boilerplate"
	}
	err := mgm.SetDefaultConfig(nil, mongodatabase, options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to Mongo at %s", mongouri)
}
