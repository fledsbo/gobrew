package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fledsbo/gobrew/graph"
	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/hwinterface"
	"github.com/fledsbo/gobrew/storage"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	monitorController := hwinterface.NewMonitorController()
	outletController := hwinterface.NewOutletController()

	storage := storage.NewStorage()
	storage.MonitorController = monitorController
	storage.OutletController = outletController

	err := storage.LoadOutlets()
	if err != nil {
		panic(err)
	}
	fermentations, err := storage.LoadFermentations()
	if err != nil {
		panic(err)
	}

	go monitorController.Scan()

	resolver := &graph.Resolver{
		MonitorController:       monitorController,
		OutletController:        outletController,
		FermentationControllers: fermentations,
		Storage:                 storage,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
