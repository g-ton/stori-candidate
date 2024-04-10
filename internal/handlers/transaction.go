package hdlGin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *HTTPHandlerGin) SendSummaryInformation(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"data": "SendSummaryInformation"})
}
