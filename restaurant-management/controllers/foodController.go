package controllers

import (
	"context"
	"net/http"
	"rtm/database"
	"rtm/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection= database.OpenCollection(database.Client,"food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")
var validate = validator.New()


func GetFoods() gin.HandlerFunc{
return func(c *gin.Context){

	}
}



func GetFood() gin.HandlerFunc{

	return  func(c *gin.Context){
		// get ctx with timeout 100
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		// get food id from parameters
		food_id := c.Param("food_id")
		// define the struct food model
		var food models.Food
		// find the ID one and decode from json to struct 
		err := foodCollection.FindOne(ctx,bson.M{"food_id":food_id}).Decode(&food)
		//put in stack for canceling db connection 
		defer cancel()
		// check err
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching food item"})
		
		}
		// send back response
		c.JSON(http.StatusOK,food)	

	}
}


func CreateFood() gin.HandlerFunc{

	return func(c *gin.Context){
				// get context and cancel 
					ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)	
				// define menu model adnd food model 

				var menu models.Menu
				var food models.Food
			// if not able to convert json -> struct -> error
				if err := c.BindJSON(&food); err!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
				return 
		}
	
	
	// validate the struct 

	validationErr := validate.Struct(food)
	
		// validation error 
			if validationErr!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
		}

		// find the menu with menu id   -> menu struct 
		err:= menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return

		}		// create food's values -> created_at, updated_at , ID , food_id , price 
		
food.Created_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
food.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

food.ID= primitive.NewObjectID()
	food.Food_id= food.ID.Hex()

var num = toFixed(*food.Price,2)
		food.Price= &num


		//try inserting 

		res,errr:= foodCollection.InsertOne(ctx,food)


		// accordingly give error or result

if errr != nil{
		msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return

		}



		defer cancel()
		c.JSON(http.StatusOK,res)

	}


}


func UpdateFood() gin.HandlerFunc{

	return func(c *gin.Context){

	}
}




func round(num float64)int{
}

func toFixed(num float64,precision int) float64{




}
