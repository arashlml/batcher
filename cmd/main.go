package main

import (
	"batcher/batcher"
	"batcher/entity"
	"batcher/handler"
	"batcher/http"
	"batcher/repository"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	connstr := "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
	db, err := repository.NewDatabase(connstr)
	if err != nil {
		log.Printf("error making the db --> %v", err)
		return
	}
	defer db.Close()
	repo := repository.NewRepository(db)

	cargo := batcher.NewBatcher[entity.Request](10, 10*time.Second, repo.BatchInsert)

	h := handler.NewHandler(cargo)

	http.RegisterAPI(e, h)

	if err := e.Start(":8080"); err != nil {
		log.Printf("Problem opening echo server %v", err.Error())
		return
	}
	log.Println("server started at port 8080")
}
