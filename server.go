package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/fledsbo/gobrew/apis"
	"github.com/fledsbo/gobrew/config"
	"github.com/fledsbo/gobrew/fermentation"
	"github.com/fledsbo/gobrew/graph"
	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/hwinterface"
	"github.com/fledsbo/gobrew/storage"
)

const defaultPort = "8080"

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	monitorController := hwinterface.NewMonitorController()
	outletController := hwinterface.NewOutletController()
	fermentationController := fermentation.NewController(monitorController, outletController)

	defer outletController.SwitchAllOff()

	storage := storage.NewStorage()
	storage.MonitorController = monitorController
	storage.OutletController = outletController
	storage.FermentationController = fermentationController

	var cfg = new(config.Config)
	config.LoadConfig(cfg)

	err := storage.LoadOutlets()
	if err != nil {
		panic(err)
	}
	err = storage.LoadFermentations()
	if err != nil {
		panic(err)
	}

	brewfather := apis.Brewfather{
		FermentationController: fermentationController,
		Config:                 cfg,
	}
	go brewfather.Run()
	go monitorController.Scan()

	resolver := &graph.Resolver{
		Config:                 cfg,
		MonitorController:      monitorController,
		OutletController:       outletController,
		FermentationController: fermentationController,
		Storage:                storage,
	}

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true // r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
