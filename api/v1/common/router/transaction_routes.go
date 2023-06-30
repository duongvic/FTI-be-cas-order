package router

import (
	apiTransaction "casorder/api/v1/common/apis/transaction"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerTransactionRoutes(r *gin.RouterGroup) {
	transactionAPI := apiTransaction.TransactionApi{}
	auth := r.Group("/transaction")
	{
		auth.GET("../transactions", transactionAPI.ListTransactions)
		auth.GET("/:id", transactionAPI.GetTransactionByID)
		auth.POST("../transactions", transactionAPI.CreateTransaction)
		auth.PATCH("/:id", transactionAPI.UpdateTransaction)
		auth.DELETE("/:id", transactionAPI.DeleteTransaction)
	}
}