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
	router.POST("/products", srv.addProduct)
	router.GET("/products/:SKU", srv.getProductWithURL)
	router.GET("/products", srv.getProductWithParam)
	router.DELETE("/products/:SKU", srv.deleteProductWithURL)
	router.DELETE("/products", srv.deleteProductWithParam)
	router.PUT("/products/:SKU", srv.updateProductWithURL)
	router.PUT("/products", srv.updateProductWithParam)
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
