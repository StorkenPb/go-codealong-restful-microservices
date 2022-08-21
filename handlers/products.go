package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/StorkenPb/restful-microservices-codealong/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l* log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// all requests will end up in here
	// Check what method is used

	// GET
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//POST (create)
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//PUT (update)
	if r.Method == http.MethodPut {
		// Expects the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		matches := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(matches) != 1 {
			p.l.Println("Invalid URI, more than one id", r.URL.Path)
			http.Error(rw, "Invalid URI - expecting product id in URI", http.StatusBadRequest)
			return
		}

		if len(matches[0]) != 2 {
			p.l.Println("path", r.URL.Path)
			p.l.Println("Invalid URI, more than two capture groups", matches)
			http.Error(rw, "Invalid URI - expecting product id in URI", http.StatusBadRequest)
			return
		}

		idString := matches[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI, unable to convert to number", idString)
			http.Error(rw, "Invalid URI - expecting product id in URI", http.StatusBadRequest)
			return
		}
		
		p.updateProduct(id, rw, r)
		return
	}

	// Catch all - If the methoud is not recognized, then a 405 is sent back
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET product")
	products := data.GetProducts()
	
	// Encode the JSON and send it to the Response writer
	err := products.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal product JSON", http.StatusInternalServerError)
	}
}

func (p*Products) addProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST product")

	product := &data.Product{}
	
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal product JSON", http.StatusBadRequest)
		return
	}

	data.AddProduc(product)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle PUT product")

	product := &data.Product{}
	
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal product JSON", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, product)
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