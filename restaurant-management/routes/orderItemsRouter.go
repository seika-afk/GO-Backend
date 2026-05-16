package routes

import (
	"github.com/gin-gonic/gin"
	controllers "rtm/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/orderItems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:order_item_id", controllers.GetOrderItem())

	incomingRoutes.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems", controllers.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:order_item_id", controllers.UpdateOrderItem())

}
