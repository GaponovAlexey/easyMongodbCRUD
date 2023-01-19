package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

type Product struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"product_name" bson:"product_name"`
	Price       int                `json:"price" bson:"price"`
	Currency    string             `json:"currency" bson:"currency"`
	Quantity    string             `json:"quantity" bson:"quantity"`
	Discount    int                `json:"discount,omitempty" bson:"discount,omitempty"`
	Vendor      string             `json:"vendor" bson:"vendor"`
	Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitempty"`
	SkuID       string             `json:"sku_id" bson:"sku_id"`
}

var (
	e   = echo.New()
	ctx context.Context
	db  *mongo.Collection
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
	db = con.Database("tronics").Collection("products")
}

func main() {

	e.GET("/", getData)
	e.POST("/", addData)

	e.Logger.Debug(e.Start(":3000"))
}

func getData(c echo.Context) error {
	var product []Product
	ca, err := db.Find(ctx, bson.M{})
	cancel(err)
	err = ca.All(ctx, &product)

	return c.JSON(http.StatusOK, product)
}

func addData(c echo.Context) error {
	var products []Product
	var insert []interface{}
	var data []interface{}

	if err := c.Bind(&products); err != nil {
		return err
	}
	for _, prod := range products {
		prod.ID = primitive.NewObjectID()
		insertId, err := db.InsertOne(ctx, prod)
		cancel(err)
		insert = append(insert, insertId.InsertedID)
		data = append(data, prod)
	}

	return c.JSON(http.StatusOK, data)
}
