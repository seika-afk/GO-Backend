package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rtm/database"
	"rtm/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		// create table and order model
		var table models.Table
		var order models.Order
		// bind json in order var -> else error badrequest
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// validate struct /our order
		// 		// handle validation error

		validateErr := validate.Struct(order)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
		}

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Item was not created"})
				return
			}

		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		result, insertErr := orderCollection.InsertOne(ctx, order)

		if insertErr != nil {
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func UpdateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

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

		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		upsert := true
		filter := bson.M{"order_id", order_id}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$st", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error", "Order item update failed"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func OrderItemOrderCreator(order models.Order) string {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id

}
