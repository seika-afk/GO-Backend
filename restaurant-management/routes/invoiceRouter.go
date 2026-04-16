package routes

import ("github.com/gin-gonic/gin"

controllers "rtm/controllers"

)

func InvoiceRoutes(incomingRoutes *gin.Engine){


incomingRoutes.GET("/invoices",controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id",controllers.GetInvoice())

incomingRoutes.POST("/invoices",controllers.CreateInvoice())

	incomingRoutes.GET("/invoices/:invoice_id",controllers.UpdateInvoice())



}
