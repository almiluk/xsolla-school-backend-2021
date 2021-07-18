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

// addProduct godoc
// @Summary add new product
// @Accept json
// @Produces json
// @Param product body models.InputProduct true "adding product"
// @Success 200 {object} models.Product "added product"
// @Failure 400 {object} models.ResponseErrorProduct
// @Failure 500 {object} models.ResponseError
// @Router /products [post]
func (srv *ProductServer) addProduct(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	newProduct := models.EmptyInputProduct()
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseErrorProduct{"json format error: " + err.Error(), *models.EmptyProduct()})
		return
	}

	if product, err := srv.db.AddProduct(*newProduct); err == nil {
		ctx.JSON(http.StatusOK, product)
	} else if err == DB.ProductAlreadyExistsError {
		ctx.JSON(http.StatusBadRequest, models.ResponseErrorProduct{err.Error(), *product})
	} else {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{err.Error()})
	}
}

// getProductWithURL godoc
// @Summary get product with specific SKU with SKU in URL path
// @Produces json
// @Param SKU path string true "SKU of searching product"
// @Success 200 {array} models.Product
// @Failure 404 {object} models.ResponseError "product with such SKU does not exist"
// @Failure 500 {object} models.ResponseError
// @Router /products/{SKU} [get]
func (srv *ProductServer) getProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	foundProduct, err := srv.db.GetProductBySKU(SKU)
	if err == nil {
		ctx.JSON(http.StatusOK, []*models.Product{foundProduct})
	} else if err == DB.ProductNotFoundError {
		ctx.JSON(http.StatusNotFound, models.ResponseError{err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{err.Error()})
	}
}

// getProductWithParam godoc
// @Summary get product with specific SKU or Id with it in URL params or all of the products, or part of them
// @Description Method return product with specific SKU, if related parameter is specified else similarly with Id.
// @Description If both of parameters aren't specified return all products or group of them, if groupSize and groupNum params are specified
// @Produces json
// @Param sku query string false "SKU of searching product"
// @Param id query int false "Id of searching product"
// @Param groupSize query int false "Size of requesting products group"
// @Param groupNum query int false "Number of requesting products group"
// @Success 200 {array} models.Product
// @Failure 404 {object} models.ResponseError "Product with specified SKU or Id not found"
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /products [get]
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
			responseData, _ = json.Marshal([]*models.Product{foundProduct})
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
			responseData, _ = json.Marshal([]*models.Product{foundProduct})
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
		ctx.JSON(code, models.ResponseError{errMsg})
	}
}

// deleteProductWithURL godoc
// @Summary delete product with specific SKU with SKU in URL path
// @Produces json
// @Param SKU path string true "SKU of deleting product"
// @Success 200
// @Failure 404 {object} models.ResponseError "product with such SKU does not exist"
// @Failure 500 {object} models.ResponseError
// @Router /products/{SKU} [delete]
func (srv *ProductServer) deleteProductWithURL(ctx *gin.Context) {
	SKU := ctx.Param("SKU")
	if err := srv.db.DeleteProductBySKU(SKU); err == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else if err == DB.ProductNotFoundError {
		ctx.JSON(http.StatusNotFound, models.ResponseError{err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{err.Error()})
	}
}

// deleteProductBySKU godoc
// @Summary delete product with specific SKU or Id with it in URL params
// @Description Method delete product with specific SKU, if related parameter is specified else similarly with Id.
// @Param sku query string false "SKU of deleting product"
// @Param id query int false "Id of deleting product"
// @Success 200
// @Failure 404 {object} models.ResponseError "Product with specified SKU or Id not found"
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /products [delete]
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
		ctx.JSON(code, models.ResponseError{errMsg})
	}
}

// updateProductBySKU godoc
// @Summary update product with specific SKU with SKU in URL path
// @Accept json
// @Produces json
// @Param product body models.InputProduct true "new product"
// @Param SKU path string false "SKU of updating product"
// @Success 200 {object} models.Product "added product"
// @Failure 400 {object} models.ResponseErrorProduct
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /products/{SKU} [PUT]
func (srv *ProductServer) updateProductWithURL(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	SKU := ctx.Param("SKU")
	newProduct := models.EmptyInputProduct()
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseErrorProduct{"json format error: " + err.Error(), *models.EmptyProduct()})
		return
	}
	if product, err := srv.db.UpdateProductBySKU(SKU, *newProduct); err == nil {
		ctx.JSON(http.StatusOK, models.ResponseErrorProduct{"", *product})
	} else if err == DB.ProductNotFoundError {
		ctx.JSON(http.StatusNotFound, models.ResponseError{err.Error()})
	} else if err == DB.ProductAlreadyExistsError {
		ctx.JSON(http.StatusBadRequest, models.ResponseErrorProduct{err.Error(), *product})
	} else {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{err.Error()})
	}
}

// updateProductWithParam godoc
// @Summary update product with specific SKU or Id with it in URL params
// @Accept json
// @Produces json
// @Param product body models.InputProduct true "new product"
// @Param sku query string false "SKU of updating product"
// @Param id query int false "Id of updating product"
// @Success 200 {object} models.Product "new product"
// @Failure 400 {object} models.ResponseErrorProduct
// @Failure 404 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /products [put]
func (srv *ProductServer) updateProductWithParam(ctx *gin.Context) {
	// TODO: More informative message about unmarshal error
	newProduct := models.EmptyInputProduct()
	err := ctx.ShouldBindJSON(newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseErrorProduct{"json format error: " + err.Error(), *models.EmptyProduct()})
		return
	}

	var errMsg string
	prod := models.EmptyProduct()
	code := http.StatusOK
	prSKU, prId, err := getSKUandIDFromUrl(ctx)
	if err != nil {
		errMsg = err.Error()
		code = http.StatusBadRequest
	} else if prSKU != "" {
		if prod, err = srv.db.UpdateProductBySKU(prSKU, *newProduct); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else if err == DB.ProductAlreadyExistsError {
				code = http.StatusBadRequest
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else if prId != 0 {
		if prod, err = srv.db.UpdateProductById(prId, *newProduct); err != nil {
			errMsg = err.Error()
			if err == DB.ProductNotFoundError {
				code = http.StatusNotFound
			} else if err == DB.ProductAlreadyExistsError {
				code = http.StatusBadRequest
			} else {
				code = http.StatusInternalServerError
			}
		}
	} else {
		errMsg = "Id or SKU of editing product must be specified"
		code = http.StatusBadRequest
	}

	if code == http.StatusOK {
		ctx.JSON(code, models.ResponseErrorProduct{"", *prod})
	} else if code == http.StatusBadRequest {
		ctx.JSON(code, models.ResponseErrorProduct{errMsg, *prod})
	} else {
		ctx.JSON(code, models.ResponseError{errMsg})
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
