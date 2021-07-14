package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	initHandlers(router)
	router.Run()
}
