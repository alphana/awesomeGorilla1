// Package handlers Products API
//
//     Schemes: http
//     Host: localhost
//     BasePath: /products
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
// swagger:meta
package handlers

import (
	"awesomeGorilla1/data"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ProductsHandler struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *ProductsHandler {
	return &ProductsHandler{l}
}

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// GetProducts ListAll handles GET requests and returns all current products
func (product *ProductsHandler) GetProducts(response http.ResponseWriter, request *http.Request) {

	product.logger.Println("GET Method Called")
	productsList := data.GetProducts()

	err := productsList.ToJson(response)

	if err != nil {
		http.Error(response, "Unable to marshal json:"+err.Error(), http.StatusInternalServerError)
	}

}

// swagger:route POST /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// PostProduct  handles post requests and returns all current products
func (product *ProductsHandler) PostProduct(response http.ResponseWriter, request *http.Request) {
	product.logger.Println("POST Method Called")
	newProductP := request.Context().Value(KeyProduct{}).(data.Product)
	data.PostProduct(&newProductP)
}

func (product *ProductsHandler) PutProduct(response http.ResponseWriter, request *http.Request) {
	product.logger.Println("PUT Method Called")

	vars := mux.Vars(request)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(response, "Unable to convert id to integer:"+err.Error(), http.StatusBadRequest)
	}

	upsertedProduct := request.Context().Value(KeyProduct{}).(data.Product)
	data.PutProduct(&upsertedProduct, ID)
}

type KeyProduct struct{}

func (product ProductsHandler) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		newProduct := data.Product{}
		err := newProduct.FromJson(request.Body)

		if err != nil {
			product.logger.Println("[ERROR] deserializing product", err)
			http.Error(responseWriter, "Error reading product", http.StatusBadRequest)
			return
		}

		err = newProduct.Validate()
		if err != nil {
			product.logger.Println("[ERROR] Validating product", err)
			http.Error(responseWriter, fmt.Sprintf("Error reading product :%s", err), http.StatusUnprocessableEntity)
			return
		}
		// Add the product to the context
		ctx := context.WithValue(request.Context(), KeyProduct{}, newProduct)
		request = request.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(responseWriter, request)
	})
}
