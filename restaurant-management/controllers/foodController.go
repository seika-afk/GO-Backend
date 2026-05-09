package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"rtm/database"
	"rtm/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//pagination -> avoiding sending too much data to frontend
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		//match stage
		matchStage := bson.D{{"$match", bson.D{{}}}}

		// group stage
		groupStage := bson.D{{
			"$group", bson.D{{
				"_id", bson.D{{
					"_id", "null"}}},

				{"total_count", bson.D{{"$sum", 1}}},

				{"data", bson.D{{"$push", "$$ROOT"}}}}}}

		// projectStage

		projectStage := bson.D{
			{"$project", bson.D{

				{"_id", 0},

				{"total_count", 1},

				{"food_items", bson.D{
					{"$slice", bson.A{"$data", startIndex, recordPerPage}}}},
			}}}

		res, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured"})
		}

		var allFoods []bson.M
		if err := res.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods[0])
	}

}

func GetFood() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get ctx with timeout 100
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// get food id from parameters
		food_id := c.Param("food_id")
		// define the struct food model
		var food models.Food
		// find the ID one and decode from json to struct
		err := foodCollection.FindOne(ctx, bson.M{"food_id": food_id}).Decode(&food)

		//put in stack for canceling db connection
		defer cancel()
		// check err
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching food item"})
			return

		}

		// send back response
		c.JSON(http.StatusOK, food)

	}
}

func CreateFood() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get context and cancel
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// define menu model adnd food model

		var menu models.Menu
		var food models.Food
		// if not able to convert json -> struct -> error
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		// validate the struct

		validationErr := validate.Struct(food)

		// validation error
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// find the menu with menu id   -> menu struct
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)

		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

			return

		} // create food's values -> created_at, updated_at , ID , food_id , price

		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()

		var num = toFixed(*food.Price, 2)
		food.Price = &num

		//try inserting

		res, errr := foodCollection.InsertOne(ctx, food)

		// accordingly give error or result

		if errr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

			return

		}

		defer cancel()
		c.JSON(http.StatusOK, res)

	}

}

func UpdateFood() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food
		var updateObj primitive.D
		food_id := c.Param("food_id")
		if err := c.BindJSON(&food); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}

		if food.Menu_id != nil {

			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Menu was not found"})

				return
			}

			updateObj = append(updateObj, bson.E{"menu_id", food.Menu_id})

		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{"food_image", food.Food_image})

		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		upsert := true
		filter := bson.M{"food_id": food_id}

		opt := options.UpdateOptions{

			Upsert: &upsert,
		}

		res, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},

			&opt,
		)
		if err != nil {
			msg := fmt.Sprint("Food item updation failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

			return
		}
		defer cancel()
		c.JSON(http.StatusAccepted, res)

	}
}

func round(num float64) int {
	return int(num * math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
