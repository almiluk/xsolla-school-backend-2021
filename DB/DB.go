package DB

import (
	"XsollaSchoolBE/models"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var ProductNotFoundError = errors.New("Product not found")
var ProductAlreadyExistsError = errors.New("Product already exists")

type DB interface {
	AddProduct(product models.InputProduct) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	GetGroupOfProducts(groupSize uint, groupNum uint) ([]*models.Product, error)
	GetProductBySKU(SKU string) (*models.Product, error)
	GetProductById(id int64) (*models.Product, error)
	DeleteProductBySKU(SKU string) error
	DeleteProductById(id int64) error
	UpdateProductBySKU(SKU string, inputProd models.InputProduct) (*models.Product, error)
	UpdateProductById(id int64, inputProd models.InputProduct) (*models.Product, error)
	Close() error
}
