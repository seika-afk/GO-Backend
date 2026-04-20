package controllers

import (
	"context"
	"log"
	"net/http"
	"rtm/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMenus() gin.HandlerFunc{


return func(c *gin.Context){
	}
	}



func GetMenu() gin.HandlerFunc{


	return func(c *gin.Context){

		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)

		res,err:= menuCollection.Find(ctx,bson.M{})
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while quering menu items"})
		}

		var allMenu[]bson.M
		if err:= res.All(ctx,&allMenu);err !=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allMenu)	


	}
}


func CreateMenu() gin.HandlerFunc{



	return func(c *gin.Context){

	}
}


func UpdateMenu() gin.HandlerFunc{

return func(c *gin.Context){

	}
}




