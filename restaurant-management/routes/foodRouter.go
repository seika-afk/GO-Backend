package routes

import ("github.com/gin-gonic/gin"
"rtm/controllers"

)






func FoodRoutes(incoming *gin.Engine){

incoming.GET("/foods",controllers.GetFoods())
incoming.GET("/foods/:food_id",controllers.GetFood())

incoming.POST("/foods",controllers.CreateFood())

	incoming.POST("/foods/:food_id",controllers.UpdateFood())

}
