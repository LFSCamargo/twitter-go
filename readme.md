<div align="center">
  <img src="https://github.com/graph-gophers/graphql-go/raw/master/docs/img/logo.png" />
  <br/>
  <h1>Go Lang GraphQL Boilerplate</h1>
  <h3>This is my GraphQL Boilerplate for starters at Go and GraphQL </h3>
</div>

# Why i created this?
I was studying Go Lang and i did not found any great tutorial of how to build a simple GraphQL server with authentication and with password hashing

# What you need to know?

There are three stuff that you need to know about this boilerplate

## What database is beign used?

MongoDB, and for managing i'm using the `mgm` that is a ODM for mongo, and the models that i create are all inside the `database/models` folder. You can create your models inside that folder like this: `database/models/user/user.go` 

## What is being used to handle GraphQL?

GQL Gen, the graphql gen it's a pretty easy solution to work with GraphQL with Go!

It has some tools to help you develop faster you just need to change the schema and run the following command to generate the models and resolvers skeleton:

```sh 
go run github.com/99designs/gqlgen
```

It will generate a model inside the file `graph/model/models_gen.go` folder and also will create the resolvers for you automatically inside `graph/schema.resolvers.go`, but there are two things that you need to know:

1. The models should not be edited because they are generated based on the `graph/schema.graphqls`
2. You can edit the `graph/schema.resolvers.go` but when you run the gql gen command will create resolvers but without any kind of logic inside

# How to run the project

1. Clone the repo
2. Enter the project folder and run the following command to install the deps `go mod download`
3. Run the `go run server.go` and you should be fine
