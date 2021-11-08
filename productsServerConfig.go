package main

import (
	"awesomeGorilla1/handlers"
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ProductServer struct {
	logger *log.Logger
}

func (serverConfig *ProductServer) ListenAndServe(logger *log.Logger) {
	serveMux := serverConfig.registerHandlers(logger)
	serverConfig.listenAndServe(logger, serveMux)
}

func (serverConfig *ProductServer) listenAndServe(logger *log.Logger, serveMux *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090" // Default port if not specified
		logger.Printf("No PORT variable found in the environment server will start using default port : %s", port)
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      serveMux,
		IdleTimeout:  120 * time.Millisecond,
		ReadTimeout:  1 * time.Millisecond,
		WriteTimeout: 1 * time.Millisecond,
	}

	go func() {
		logger.Printf("Starting server on : %s", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Kill)
	signalInstance := <-signalChannel
	logger.Println("Received terminate, will shutdown", signalInstance)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

func (serverConfig *ProductServer) registerHandlers(logger *log.Logger) *mux.Router {
	productHandler := handlers.NewProducts(logger)

	serveMux := mux.NewRouter()
	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)
	postRouter := serveMux.Methods("POST").Subrouter()
	postRouter.HandleFunc("/products", productHandler.PostProduct)
	postRouter.Use(productHandler.MiddlewareValidateProduct)
	putRouter := serveMux.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.PutProduct)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	redocOptions := middleware.RedocOpts{SpecURL: "./swagger.json"}
	redoc := middleware.Redoc(redocOptions, nil)
	getRouter.Handle("/docs", redoc)
	getRouter.Handle("/swagger.json", http.FileServer(http.Dir("./")))

	return serveMux
}
