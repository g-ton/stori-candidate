package websrv

import (
	"github.com/g-ton/stori-candidate/internal/core/ports"
)

type webService struct {
	web ports.WebService
}

func New() *webService {
	return &webService{}
}

func (webSrv *webService) Run() {}
