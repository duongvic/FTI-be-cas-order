package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type TransactionApi struct {
	services.CommonApi
}

// ListTransactions Gets the list of all transactions
func (t TransactionApi) ListTransactions(c *gin.Context) {
	if err := t.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		t.Error(500, err, err.Error())
		t.Logger.Error(err.Error())
		return
	}
	
	var transaction models.Transaction
	var transactions []*models.Transaction
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := t.List(transaction, &transactions,&pagination, t.Orm)
	if err != nil {
		t.Logger.Errorf("Error getting transactions: %v", err)
		t.Error(404, err, "Failed to get transactions")
		return
	}

	t.OK(result, "Success")
}

// GetTransactionByID Gets specific transaction from ID in context
func (t TransactionApi) GetTransactionByID(c *gin.Context) {
	if err := t.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		t.Error(500, err, err.Error())
		t.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var transaction models.Transaction

	result, err := t.Find(&transaction, id, t.Orm)
	if err != nil {
		t.Logger.Errorf("Error getting transaction: %v", err)
		t.Error(404, err, "Failed to get transaction")
		return
	}

	t.OK(result, "Success")
}

// CreateTransaction Creates a new transaction
func (t TransactionApi) CreateTransaction(c *gin.Context) {
	if err := t.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		t.Error(500, err, err.Error())
		t.Logger.Error(err.Error())
		return
	}
	
	var model models.Transaction

	transaction, err := t.ParseObjectFromRequest(c, &model)
	if err != nil {
		t.Logger.Errorf("Error creating transaction: %v", err)
		t.Error(400, err, "Failed to create transaction")
	}

	if err := t.Orm.Create(transaction).Error; err != nil {
		t.Logger.Errorf("Error creating transaction: %v", err)
		t.Error(400, err, "Failed to create transaction")
		return
	}

	t.OK(transaction, "Success")
}

// UpdateTransaction Updates a transaction
func (t TransactionApi) UpdateTransaction(c *gin.Context) {
	if err := t.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		t.Error(500, err, err.Error())
		t.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.Transaction

	transaction, err := t.ParseObjectFromRequest(c, &model)
	if err != nil {
		t.Logger.Errorf("Error reading data: %v", err)
		t.Error(400, err, "Failed to update transaction")
		return
	}

	if err := t.Orm.Model(&model).Where("id = ?", id).Updates(transaction).Error; err != nil {
		t.Logger.Errorf("Error creating transaction: %v", err)
		t.Error(400, err, "Failed to update transaction")
		return
	}
	result, err := t.Find(&model, id, t.Orm)
	t.OK(result, "Success")
}

// DeleteTransaction Deletes a transaction
func (t TransactionApi) DeleteTransaction(c *gin.Context) {
	if err := t.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		t.Error(500, err, err.Error())
		t.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var transaction models.Transaction

	if err := t.Orm.First(&transaction, "id = ?", id).Delete(&transaction).Error; err != nil {
		t.Logger.Errorf("Error deleting transaction: %v", err)
		t.Error(404, err, "Failed to delete transaction")
		return
	}
	t.OK(transaction, "Success")
}