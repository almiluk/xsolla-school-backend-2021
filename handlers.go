package main

import (
	"XollaSchoolBE/DB"
	"XollaSchoolBE/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

var db *DB.DB

func initHandlers(router *gin.Engine) {
	var err error
	db, err = DB.Init("testDB.db")
	if err != nil {
		panic(err)
	}
	router.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"Status": "It is working"}) })
	router.POST("/products", addProduct)
	router.GET("/products/:SKU", getProductWithURL)
	router.GET("/products", getProductWithParam)
	router.DELETE("/products/:SKU", deleteProductWithURL)
	router.DELETE("/products", deleteProductWithParam)
	router.PUT("/products/:SKU", updateProductWithURL)
	router.PUT("/products", updateProductWithParam)
}

// TODO: Check if request body signature is matches with models.InputProduct

func addProduct(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	newProduct := models.EmptyInputProduct()
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "json format error: " + err.Error()})
		return
	} else if newProduct.Name == "" || newProduct.SKU == "" || newProduct.Type == "" {

	}

	if product, err := db.AddProduct(*newProduct); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "product": product})
	}
}

func getProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	foundProduct, err := db.GetProductBySKU(SKU)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
	} else {
		ctx.JSON(http.StatusOK, foundProduct)
	}
}

func getProductWithParam(ctx *gin.Context) {
	var jsonData []byte
	var errMsg string
	code := http.StatusOK
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		if foundProduct, err := db.GetProductBySKU(prSKU); err == nil {
			jsonData, _ = json.Marshal(foundProduct)
		} else {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prIdStr, ok := ctx.GetQuery("id"); ok {
		prId, err := strconv.ParseInt(prIdStr, 10, 64)
		if err != nil {
			errMsg = err.Error()
			code = http.StatusBadRequest
		} else {
			if foundProduct, err := db.GetProductByID(prId); err == nil {
				jsonData, _ = json.Marshal(foundProduct)
			} else {
				errMsg = err.Error()
				if err == DB.ProductNotFoundError {
					code = http.StatusNotFound
				} else {
					code = http.StatusInternalServerError
				}
			}

		}
	} else {
		products, err := db.GetAllProducts()
		if err != nil {
			errMsg = err.Error()
			code = http.StatusInternalServerError
		} else {
			jsonData, _ = json.Marshal(products)
		}
	}

	if errMsg == "" {
		ctx.Header("Content-Type", "application/json")
		ctx.String(code, string(jsonData))
	} else {
		ctx.JSON(code, gin.H{"error": errMsg})
	}
}

func deleteProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	if err := db.DeleteProductBySKU(SKU); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
	}
}

func deleteProductWithParam(ctx *gin.Context) {
	var errMsg string
	code := http.StatusOK
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		if err := db.DeleteProductBySKU(prSKU); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prIdStr, ok := ctx.GetQuery("id"); ok {
		prId, err := strconv.ParseInt(prIdStr, 10, 64)
		if err != nil {
			errMsg = err.Error()
			code = http.StatusBadRequest
		} else if err := db.DeleteProductByID(prId); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else {
		errMsg = "Id or SKU of deleting product must be specified"
		code = http.StatusBadRequest
	}

	if errMsg == "" {
		ctx.String(code, "")
	} else {
		ctx.JSON(code, gin.H{"error": errMsg})
	}
}

func updateProductWithURL(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	SKU := ctx.Param("SKU")
	newProduct := models.EmptyInputProduct()
	bodyData, _ := ioutil.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(bodyData, newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := db.UpdateProductBySKU(SKU, *newProduct); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else if err == DB.ProductNotFoundError {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func updateProductWithParam(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	newProduct := models.EmptyInputProduct()
	bodyData, _ := ioutil.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(bodyData, newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var errMsg string
	code := http.StatusOK
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		if err := db.UpdateProductBySKU(prSKU, *newProduct); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prIdStr, ok := ctx.GetQuery("id"); ok {
		prId, err := strconv.ParseInt(prIdStr, 10, 64)
		if err != nil {
			errMsg = err.Error()
			code = http.StatusBadRequest
		} else if err := db.UpdateProductById(prId, *newProduct); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else {
		errMsg = "Id or SKU of editing product must be specified"
		code = http.StatusBadRequest
	}

	if errMsg == "" {
		ctx.String(code, "")
	} else {
		ctx.JSON(code, gin.H{"error": errMsg})
	}
}
