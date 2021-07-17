package productServer

import (
	"XollaSchoolBE/DB"
	"XollaSchoolBE/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// TODO: Check if request body signature is matches with models.InputProduct

func (srv *ProductServer) addProduct(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	newProduct := models.EmptyInputProduct()
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "json format error: " + err.Error()})
		return
	} else if newProduct.Name == "" || newProduct.SKU == "" || newProduct.Type == "" {

	}

	if product, err := srv.db.AddProduct(*newProduct); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "product": product})
	}
}

func (srv *ProductServer) getProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	foundProduct, err := srv.db.GetProductBySKU(SKU)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, foundProduct)
	}
}

func (srv *ProductServer) getProductWithParam(ctx *gin.Context) {
	var jsonData []byte
	var errMsg string
	code := http.StatusOK
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		if foundProduct, err := srv.db.GetProductBySKU(prSKU); err == nil {
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
			if foundProduct, err := srv.db.GetProductByID(prId); err == nil {
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
		products, err := srv.db.GetAllProducts()
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

func (srv *ProductServer) deleteProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	if err := srv.db.DeleteProductBySKU(SKU); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func (srv *ProductServer) deleteProductWithParam(ctx *gin.Context) {
	var errMsg string
	code := http.StatusOK
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		if err := srv.db.DeleteProductBySKU(prSKU); err != nil {
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
		} else if err := srv.db.DeleteProductByID(prId); err != nil {
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

func (srv *ProductServer) updateProductWithURL(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	SKU := ctx.Param("SKU")
	newProduct := models.EmptyInputProduct()
	bodyData, _ := ioutil.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(bodyData, newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := srv.db.UpdateProductBySKU(SKU, *newProduct); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else if err == DB.ProductNotFoundError {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if err == DB.ProductAlreadyExistsError {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (srv *ProductServer) updateProductWithParam(ctx *gin.Context) {
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
		if err := srv.db.UpdateProductBySKU(prSKU, *newProduct); err != nil {
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
		} else if err := srv.db.UpdateProductById(prId, *newProduct); err != nil {
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
