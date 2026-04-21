package controllers

import (
	"context"
	"log"
	"net/http"
	"rtm/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMenus() gin.HandlerFunc{


return func(c *gin.Context){
// get context and cancel
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		// get menu id 
		menu_id := c.Param("menu_id")
		//create a menu model 
		var menu models.Menu
		// try finding the foodCollection with that menu id and put in menu model 
		err:= foodCollection.FindOne(ctx,bson.M{"menu_id":menu_id}).Decode(&menu)

	
		// put cancel in stack 
	defer cancel()
		// err 
		if err != nil{

			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching menu"})
		}
		//return 
c.JSON(http.StatusOK,menu)
		
		
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

var menu models.Menu
if err:= c.BindJSON(&menu);err != nil{
	c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	return
}

ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
defer cancel()

menu.Created_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
menu.Updated_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
menu.ID = primitive.NewObjectID()
menu.Menu_id = menu.ID.Hex()

result,err := menuCollection.InsertOne(ctx,menu)

if err != nil{
	c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while creating menu item"})
	return
}
c.JSON(http.StatusOK,result)	
	
	}
}


func UpdateMenu() gin.HandlerFunc{

return func(c *gin.Context){

	}
}




