package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/khanhct/go-lib-core/sdk/api"

	"casorder/db/models"
)

type HealthApi struct {
	api.Api
}

// check Health of Service

func (h HealthApi) Check(c *gin.Context) {
	if err := h.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		h.Error(500, err, err.Error())
		h.Logger.Error(err.Error())
		return
	}
	object := models.Health{
		StatusCode: 200,
		Message:    "I'm OK",
	}
	h.OK(object, "Success")
}
