package DB

import (
	"XollaSchoolBE/models"
	"errors"
)

var ProductNotFoundError = errors.New("Product not found")
var ProductAlreadyExistsError = errors.New("Product already exists")

var products []models.Product

func AddProduct(prod models.InputProduct) error {
	if _, err := GetProductBySKU(prod.SKU); err == nil {
		return ProductAlreadyExistsError
	} else {
		products = append(products, *models.NewProduct(prod.SKU, prod.Name, prod.Type, prod.Cost, uint64(len(products))))
		return nil
	}
}

func GetAllProducts() []models.Product {
	return products
}

func GetProductByID(id uint64) (models.Product, error) {
	for _, product := range products {
		if product.Id == id {
			return product, nil
		}
	}
	return *models.EmptyProduct(), ProductNotFoundError
}

func GetProductBySKU(SKU string) (models.Product, error) {
	for _, product := range products {
		if product.SKU == SKU {
			return product, nil
		}
	}
	return *models.EmptyProduct(), ProductNotFoundError
}

func DeleteProductByID(id uint64) error {
	for i, product := range products {
		if product.Id == id {
			products = append(products[:i], products[i+1])
			return nil
		}
	}
	return ProductNotFoundError
}

func DeleteProductBySKU(SKU string) error {
	for i, product := range products {
		if product.SKU == SKU {
			products = append(products[:i], products[i+1:]...)
			return nil
		}
	}
	return ProductNotFoundError
}

func ReplaceProductBySKU(oldSKU string, inputProd models.InputProduct) error {
	for _, prod := range products {
		if prod.SKU == oldSKU {
			prod.Name = inputProd.Name
			prod.Type = inputProd.Type
			prod.Cost = inputProd.Cost
			return nil
		}
	}
	return ProductNotFoundError
}

func ReplaceProductByID(ID uint64, inputProd models.InputProduct) error {
	for _, prod := range products {
		if prod.Id == ID {
			prod.Name = inputProd.Name
			prod.Type = inputProd.Type
			prod.Cost = inputProd.Cost
			return nil
		}
	}
	return ProductNotFoundError
}
