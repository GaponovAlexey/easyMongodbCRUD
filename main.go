package main

import (
	"context"
	"encoding/json"
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
	DB  *mongo.Collection
)

func cancel(e error) {
	if e != nil {
		fmt.Errorf("This Error <---- ", e)
	}
}

func init() {
	url := "mongodb://localhost:27017"
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	cancel(err)
	DB = con.Database("tronics").Collection("products")
}

func main() {

	e.GET("/", getData)
	e.GET("/:id", getDataID)
	e.POST("/", addData)
	e.PUT("/:id", putData)
	e.DELETE("/:id", deleteProd)

	e.Logger.Debug(e.Start(":3000"))
}

// get
func getData(c echo.Context) error {
	var prod []Product

	ca, err := DB.Find(ctx, bson.M{})
	cancel(err)
	ca.All(ctx, &prod)
	return c.JSON(http.StatusOK, prod)
}

// getId
func getDataID(c echo.Context) error {
	var prodId Product
	gId, err := primitive.ObjectIDFromHex(c.Param("id"))
	cancel(err)
	filter := bson.M{"_id": gId}

	f := DB.FindOne(ctx, filter)

	if err := f.Decode(&prodId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, prodId)
}

//POST 
func addData(c echo.Context) error {
	var products Product
	if err := c.Bind(&products); err != nil {
		return err
	}
	// for _, prod := range products {

	products.ID = primitive.NewObjectID()
	_, err := DB.InsertOne(ctx, products)
	cancel(err)
	// }
	return c.JSON(http.StatusOK, "success")
}
//PUT PATCH
func putData(c echo.Context) error {
	var prod Product

	gID, err := primitive.ObjectIDFromHex(c.Param("id"))

	cancel(err)
	filter := bson.M{"_id": gID}
	f := DB.FindOne(ctx, filter)

	if err := f.Decode(&prod); err != nil {
		return err
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&prod); err != nil {
		return err
	}

	newData := bson.M{"$set": prod}

	DB.UpdateOne(ctx, filter, newData)

	return c.JSON(http.StatusOK, prod)
}
//DELETE
func deleteProd(c echo.Context) error {
	gID, err := primitive.ObjectIDFromHex(c.Param("id"))
	cancel(err)
	filter := bson.M{"_id": gID}

	_, err = DB.DeleteOne(ctx, filter)
	cancel(err)

	return c.JSON(http.StatusOK, "success")
}
