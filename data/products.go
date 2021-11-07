package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"regexp"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku" validate:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string
}

func (product *Product) Validate() error {
	validatorr := validator.New()
	validatorr.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		regex := regexp.MustCompile("[a-zA-Z0-9_]+") // \w+ or  [\w-]+
		match := regex.FindAllString(fl.Field().String(), -1)

		return len(match) == 1
	})
	return validatorr.Struct(product)
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func PostProduct(product *Product) {
	product.ID = getNextID()
	productList = append(productList, product)
}

func PutProduct(product *Product, ID int) error {
	_, index, err := GetProduct(ID)
	if err != nil {
		if err == ErrProductNotFound {
			PostProduct(product)
			return nil
		} else {
			return err
		}

	}
	product.ID = ID
	productList[index] = product
	return nil
}

func getNextID() int {
	return productList[len(productList)-1].ID + 1
}

var ErrProductNotFound = fmt.Errorf("product not found")

func GetProduct(ID int) (*Product, int, error) {
	for i, product := range productList {
		if product.ID == ID {
			return product, i, nil
		}

	}
	return nil, -1, ErrProductNotFound
}

func (product *Product) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(product)
}

func (product *Products) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(product)

}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
