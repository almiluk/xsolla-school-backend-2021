package productServer

import (
	"XsollaSchoolBE/DB"
	"XsollaSchoolBE/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
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
	}

	if product, err := srv.db.AddProduct(*newProduct); err == nil {
		ctx.JSON(http.StatusOK, product)
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
	var responseData []byte
	var errMsg string
	code := http.StatusOK
	prSKU, prId, err := getSKUandIDFromUrl(ctx)
	if err != nil {
		errMsg = err.Error()
		code = http.StatusBadRequest
	} else if prSKU != "" {
		if foundProduct, err := srv.db.GetProductBySKU(prSKU); err == nil {
			responseData, _ = json.Marshal(foundProduct)
		} else {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prId != 0 {
		if foundProduct, err := srv.db.GetProductByID(prId); err == nil {
			responseData, _ = json.Marshal(foundProduct)
		} else {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else {
		var products []models.Product
		var err error
		groupSizeStr, okSize := ctx.GetQuery("groupSize")
		groupNumStr, okNum := ctx.GetQuery("groupNum")
		if okSize && okNum {
			if groupSize, err := strconv.ParseUint(groupSizeStr, 10, 32); err != nil {
				errMsg = "groupSize parameter must be an 32-bit unsigned integer"
			} else if groupNum, err := strconv.ParseUint(groupNumStr, 10, 32); err != nil {
				errMsg = "groupNum parameter must be an 32-bit unsigned integer"
			} else {
				products, err = srv.db.GetGroupOfProducts(uint(groupSize), uint(groupNum))
			}
		} else {
			products, err = srv.db.GetAllProducts()
		}
		if err != nil {
			errMsg = err.Error()
			code = http.StatusInternalServerError
		} else {
			responseData, _ = json.Marshal(products)
		}
	}

	if errMsg == "" {
		ctx.Header("Content-Type", "application/json")
		ctx.String(code, string(responseData))
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
	prSKU, prId, err := getSKUandIDFromUrl(ctx)
	if err != nil {
		errMsg = err.Error()
		code = http.StatusBadRequest
	} else if prSKU != "" {
		if err := srv.db.DeleteProductBySKU(prSKU); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prId != 0 {
		if err := srv.db.DeleteProductByID(prId); err != nil {
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
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "json format error: " + err.Error()})
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
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "json format error: " + err.Error()})
		return
	}

	var errMsg string
	code := http.StatusOK
	prSKU, prId, err := getSKUandIDFromUrl(ctx)
	if err != nil {
		errMsg = err.Error()
		code = http.StatusBadRequest
	} else if prSKU != "" {
		if err := srv.db.UpdateProductBySKU(prSKU, *newProduct); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prId != 0 {
		if err := srv.db.UpdateProductById(prId, *newProduct); err != nil {
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

func getSKUandIDFromUrl(ctx *gin.Context) (string, int64, error) {
	if prSKU, ok := ctx.GetQuery("sku"); ok {
		return prSKU, 0, nil
	} else if prIdStr, ok := ctx.GetQuery("id"); ok {
		if prId, err := strconv.ParseInt(prIdStr, 10, 64); err != nil {
			return "", 0, errors.New("id parameter type error:" + err.Error())
		} else {
			return "", prId, nil
		}
	} else {
		return "", 0, nil
	}
}
