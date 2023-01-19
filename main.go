package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var (
	e   = echo.New()
	ctx context.Context
)

func cancel(e error) {
	if e != nil {
		fmt.Errorf("This Error <---- ", e)
	}
}

func init() {
	url := fmt.Sprintf("mongodb://%s:%s", "localhost", "27017")
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	cancel(err)
	col := con.Database("tronics")
	db := col.Collection("products")
}

func main() {

	e.GET("/", getData)

	e.Logger.Debug(e.Start(":3000"))
}

func getData(c echo.Context) error {

	return c.JSON(http.StatusOK, "ok")
}
