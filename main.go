//go:generate swagger generate spec
package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "rest-api ", log.LstdFlags)
	productServer := ProductServer{logger: logger}
	productServer.ListenAndServe(logger)
}
