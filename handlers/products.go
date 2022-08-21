package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/StorkenPb/restful-microservices-codealong/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l* log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET product")
	products := data.GetProducts()
	
	// Encode the JSON and send it to the Response writer
	err := products.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal product JSON", http.StatusInternalServerError)
	}
}

func (p*Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST product", r.Context().Value(KeyProduct{}))

	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduc(&product)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle PUT product")

	// Id comes from the router as a map
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id to integer", http.StatusBadRequest)
	}

	// product comes from the middleware
	product := r.Context().Value(KeyProduct{}).(data.Product)
	
	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFount {
		p.l.Println(err)
		http.Error(rw, "Product not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

	rw.Write([]byte("Successfully updated product"))
}

type KeyProduct struct {}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}
		
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// Passing the product as a value in the context 
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}