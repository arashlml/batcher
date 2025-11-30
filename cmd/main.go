package main

import (
	batcherV3 "batcher/batcher3"
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
	connStr := "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
	db, err := repository.NewDatabase(connStr)
	if err != nil {
		log.Printf("error making the db --> %v", err)
		return
	}
	defer db.Close()
	repo := repository.NewRepository(db)

	cargo := batcherV3.NewBatcher[entity.Request](10, 10*time.Second)

	h := handler.NewHandler(cargo, repo)

	http.RegisterAPI(e, h)

	if err := e.Start(":8080"); err != nil {
		log.Printf("Problem opening echo server %v", err.Error())
		return
	}
	log.Println("server started at port 8080")
}
