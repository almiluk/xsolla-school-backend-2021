package DB

import (
	"XsollaSchoolBE/models"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

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
	"getGroupOfProducts": "SELECT * FROM Products ORDER BY id LIMIT ? OFFSET ?",
	"insertProduct":      "INSERT INTO Products(SKU, name, type, cost) VALUES(?, ?, ?, ?)",
	"deleteProductBySKU": "DELETE FROM Products WHERE SKU=?",
	"deleteProductById":  "DELETE FROM Products WHERE id=?",
	"updateProductBySKU": "UPDATE Products SET SKU=?, name=?, type=?, cost=? WHERE SKU=?",
	"updateProductById":  "UPDATE Products SET SKU=?, name=?, type=?, cost=? WHERE id=?",
}

type sqlite3DB struct {
	*sql.DB
}

func InitSqlite3DB(DBfilename string) (*sqlite3DB, error) {
	var err error
	sqlDB, err := sql.Open("sqlite3", DBfilename)
	if err != nil {
		return nil, fmt.Errorf("db init error: %v", err)
	}
	db := sqlite3DB{sqlDB}
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

func (db *sqlite3DB) AddProduct(product models.InputProduct) (*models.Product, error) {
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

func (db *sqlite3DB) GetAllProducts() ([]*models.Product, error) {
	rows, err := db.Query(sqlQueries["getAllProducts"])
	if err != nil {
		return nil, err
	}
	products := make([]*models.Product, 0)
	for rows.Next() {
		var product models.Product
		rows.Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
		products = append(products, &product)
	}
	return products, nil
}

func (db *sqlite3DB) GetGroupOfProducts(groupSize uint, groupNum uint) ([]*models.Product, error) {
	rows, err := db.Query(sqlQueries["getGroupOfProducts"], groupSize, (groupNum-1)*groupSize)
	if err != nil {
		return nil, err
	}
	products := make([]*models.Product, 0)
	for i := 0; rows.Next(); i++ {
		var product models.Product
		rows.Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
		products = append(products, &product)
	}
	return products, nil
}

func (db *sqlite3DB) GetProductBySKU(SKU string) (*models.Product, error) {
	var product models.Product
	err := db.QueryRow(sqlQueries["getProductBySKU"], SKU).Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
	if err == sql.ErrNoRows {
		return nil, ProductNotFoundError
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func (db *sqlite3DB) GetProductById(id int64) (*models.Product, error) {
	var product models.Product
	err := db.QueryRow(sqlQueries["getProductById"], id).Scan(&product.Id, &product.SKU, &product.Name, &product.Type, &product.Cost)
	if err == sql.ErrNoRows {
		return nil, ProductNotFoundError
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func (db *sqlite3DB) DeleteProductById(id int64) error {
	_, err := db.GetProductById(id)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["deleteProductById"], id)
	return err
}

func (db *sqlite3DB) DeleteProductBySKU(SKU string) error {
	_, err := db.GetProductBySKU(SKU)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlQueries["deleteProductBySKU"], SKU)
	return err
}

func (db *sqlite3DB) UpdateProductBySKU(SKU string, inputProd models.InputProduct) (*models.Product, error) {
	prod, err := db.GetProductBySKU(SKU)
	if err != nil {
		return prod, err
	}
	_, err = db.Exec(sqlQueries["updateProductBySKU"], inputProd.SKU, inputProd.Name, inputProd.Type, inputProd.Cost, SKU)
	if sqlErr, ok := err.(sqlite3.Error); ok && sqlErr.ExtendedCode == 2067 {
		// Error code 2067 means UNIQUE constraint failed (https://www.sqlite.org/rescode.html#constraint_unique)
		prod, _ = db.GetProductBySKU(inputProd.SKU)
		return prod, ProductAlreadyExistsError
	}
	return &models.Product{inputProd, prod.Id}, err
}

func (db *sqlite3DB) UpdateProductById(id int64, inputProd models.InputProduct) (*models.Product, error) {
	prod, err := db.GetProductById(id)
	if err != nil {
		return prod, err
	}
	_, err = db.Exec(sqlQueries["updateProductById"], inputProd.SKU, inputProd.Name, inputProd.Type, inputProd.Cost, id)
	if sqlErr, ok := err.(sqlite3.Error); ok && sqlErr.ExtendedCode == 2067 {
		// Error code 2067 means UNIQUE constraint failed (https://www.sqlite.org/rescode.html#constraint_unique)
		prod, _ = db.GetProductBySKU(inputProd.SKU)
		return prod, ProductAlreadyExistsError
	}
	return &models.Product{inputProd, prod.Id}, err
}
