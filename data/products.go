package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string `json:"sku" validate:"required,sku"`
	CreatedAT   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-asdf-asdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func GetProducts() Products {
	return productList
}

func UpdateProduct(id int, p *Product) error {
	pos, err := findProductPosition(id)
	if err != nil {
		return err
	}

	// id is taken from the url path, p.ID can be something else...
	p.ID = id
	productList[pos] = p

	return nil
}

func AddProduc(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	lp := productList[len(productList) -1]
	return lp.ID + 1
}

var ErrProductNotFount = fmt.Errorf("Product not found")

func findProductPosition(id int) (int, error) {
	for i, p := range productList {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrProductNotFount
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedAT:   time.Now().Local().String(),
		UpdatedAt:   time.Now().Local().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "def456",
		CreatedAT:   time.Now().Local().String(),
		UpdatedAt:   time.Now().Local().String(),
	},
}