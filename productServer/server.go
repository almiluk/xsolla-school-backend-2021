package productServer

import (
	"XsollaSchoolBE/DB"
	_ "XsollaSchoolBE/docs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type ProductServer struct {
	*http.Server
	db DB.DB
}

func Run(addr string, DBfilename string) (*ProductServer, error) {
	var err error
	db, err := DB.InitSqlite3DB(DBfilename)
	if err != nil {
		return nil, err
	}
	srv := ProductServer{&http.Server{Addr: addr}, db}
	srv.initHandlers()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	return &srv, nil
}

func (srv *ProductServer) initHandlers() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"Status": "It is working"}) })
	v1ProductsGroup := router.Group("api/v1/products")
	{
		v1ProductsGroup.POST("", srv.addProduct)
		v1ProductsGroup.GET("/:SKU", srv.getProductWithURL)
		v1ProductsGroup.GET("", srv.getProductWithParam)
		v1ProductsGroup.HEAD("/:SKU", srv.headProductsWithURL)
		v1ProductsGroup.HEAD("", srv.headProductsWithParam)
		v1ProductsGroup.DELETE("/:SKU", srv.deleteProductWithURL)
		v1ProductsGroup.DELETE("", srv.deleteProductWithParam)
		v1ProductsGroup.PUT("/:SKU", srv.updateProductWithURL)
		v1ProductsGroup.PUT("", srv.updateProductWithParam)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	srv.Handler = router
}

func (srv *ProductServer) Shutdown(ctx context.Context) error {
	servErr := srv.Server.Shutdown(ctx)
	DBErr := srv.db.Close()
	if servErr != nil {
		return servErr
	} else {
		return DBErr
	}
}
