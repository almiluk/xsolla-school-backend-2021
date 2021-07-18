package main

import (
	"XollaSchoolBE/productServer"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	srv, err := productServer.Run(":8080", "products.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
