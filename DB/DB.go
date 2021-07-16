package DB

import (
	"XollaSchoolBE/models"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var ProductNotFoundError = errors.New("Product not found")
var ProductAlreadyExistsError = errors.New("Product already exists")

var sqlQueries = map[string]string{
	"init": `
	CREATE TABLE IF NOT EXISTS Products (
		id INTEGER PRIMARY KEY,
		SKU TEXT,
		name TEXT,
		type TEXT,
		cost INTEGER,
		UNIQUE(SKU)
	)`,
	"getProductById":     "SELECT * FROM Products WHERE id=?",
	"getProductBySKU":    "SELECT * From Products WHERE SKU=?",
	"getAllProducts":     "SELECT * FROM Products",
	"insertProduct":      "INSERT INTO Products(SKU, name, type, cost) VALUES(?, ?, ?, ?)",
	"deleteProductBySKU": "DELETE FROM Products WHERE SKU=?",
	"deleteProductById":  "DELETE FROM Products WHERE id=?",
	"updateProductBySKU": "UPDATE Products SET SKU=?, name=?, type=?, cost=? WHERE SKU=?",
	"updateProductById":  "UPDATE Products SET SKU=?, name=?, type=?, cost=? WHERE id=?",
}

type DB struct {
	*sql.DB
}

func Init(DBfilename string) (*DB, error) {
	var err error
	sqlDB, err := sql.Open("sqlite3", DBfilename)
	if err != nil {
		return nil, fmt.Errorf("db init error: %v", err)
	}
	db := DB{sqlDB}
	_, err = db.Exec(sqlQueries["init"])
	if err != nil {
		return nil, err
	}
	for _, query := range sqlQueries {
		_, err = db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("sql query prepare error: %v\n%s", err, query)
		}
	}
	return &db, nil
}

func (db *DB) AddProduct(product models.InputProduct) (*models.Product, error) {
	prod, err := db.GetProductBySKU(product.SKU)
	if err == ProductNotFoundError {
		res, err := db.Exec(sqlQueries["insertProduct"], product.SKU, product.Name, product.Type, product.Cost)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		return &models.Product{product, id}, nil
	} else if err != nil {
		return nil, err
	} else {
		return prod, ProductAlreadyExistsError
	}
}

func (db *DB) GetAllProducts() ([]models.Product, error) {
	rows, err := db.Query(sqlQueries["getAllProducts"])
	if err != nil {
		return nil, err
	}
	products := make([]models.Product, 0)
	var product models.Product
	for rows.Next() {
		rows.Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
		products = append(products, product)
	}
	return products, nil
}

func (db *DB) GetProductBySKU(SKU string) (*models.Product, error) {
	var product models.Product
	err := db.QueryRow(sqlQueries["getProductBySKU"], SKU).Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
	if err == sql.ErrNoRows {
		return nil, ProductNotFoundError
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func (db *DB) GetProductByID(id int64) (*models.Product, error) {
	var product models.Product
	err := db.QueryRow(sqlQueries["getProductById"], id).Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
	if err == sql.ErrNoRows {
		return nil, ProductNotFoundError
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func (db *DB) DeleteProductByID(id int64) error {
	_, err := db.GetProductByID(id)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["deleteProductById"], id)
	return err
}

func (db *DB) DeleteProductBySKU(SKU string) error {
	_, err := db.GetProductBySKU(SKU)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["deleteProductBySKU"], SKU)
	return err
}

func (db *DB) UpdateProductBySKU(SKU string, inputProd models.InputProduct) error {
	_, err := db.GetProductBySKU(SKU)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["updateProductBySKU"], inputProd.SKU, inputProd.Name, inputProd.Type, inputProd.Cost, SKU)
	return err
}

func (db *DB) UpdateProductById(id int64, inputProd models.InputProduct) error {
	_, err := db.GetProductByID(id)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["updateProductById"], inputProd.SKU, inputProd.Name, inputProd.Type, inputProd.Cost, id)
	return err
}

/*var products []models.Product

func AddProduct(prod models.InputProduct) error {
	if _, err := GetProductBySKU(prod.SKU); err == nil {
		return ProductAlreadyExistsError
	} else {
		products = append(products, *models.NewProduct(prod.SKU, prod.Name, prod.Type, prod.Cost, int64(len(products))))
		return nil
	}
}

func GetAllProducts() []models.Product {
	return products
}

func GetProductByID(id int64) (models.Product, error) {
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

func DeleteProductByID(id int64) error {
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

func ReplaceProductByID(ID int64, inputProd models.InputProduct) error {
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

*/
