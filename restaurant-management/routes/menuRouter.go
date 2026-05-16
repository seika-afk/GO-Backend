package routes

import (
	"github.com/gin-gonic/gin"

	controllers "rtm/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/menus", controllers.GetMenus())
	incomingRoutes.GET("/menus/:menu_id", controllers.GetMenus())

	incomingRoutes.POST("/menus", controllers.CreateMenu())

	incomingRoutes.PATCH("/menus/:menu_id", controllers.UpdateMenu())

}
