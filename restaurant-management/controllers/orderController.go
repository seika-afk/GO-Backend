package controllers

import (
	"context"
	"log"
	"net/http"
	"rtm/database"
	"rtm/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing other items"})

		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		// get ctx with timeout 100
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// get food id from parameters
		order_id := c.Param("order_id")
		// define the struct food model
		var order models.Order
		// find the ID one and decode from json to struct
		err := foodCollection.FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order)
		defer cancel()
		// check err
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching food item"})
			return

		}

		// send back response
		c.JSON(http.StatusOK, order)

	}
}

func CreateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		var updateObj primitive.D
		order_id := c.Param("order_id")

		if order.Table_id != nil {

			err := menuCollection.FindOne(ctx, bson.M{"table_id": food.Table_id}).Decode(&table)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Menu was not found"})

				return
			}

			updateObj = append(updateObj, bson.E{"menu", order.Table_id})

		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	}
}
