package hdlGin

import (
	"github.com/g-ton/stori-candidate/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// Here We are implementing "Dependency injection" connecting Ports with Adapters

type HTTPHandlerGin struct {
	dbService  ports.DatabaseService
	Router     *gin.Engine
	webService ports.WebService
}

func NewHTTPHandlerGin(dbService ports.DatabaseService, webService ports.WebService) *HTTPHandlerGin {
	return &HTTPHandlerGin{
		dbService:  dbService,
		Router:     gin.New(),
		webService: webService,
	}
}

func (hdl *HTTPHandlerGin) Run() {
	hdl.Router.Run(":9696")
}
