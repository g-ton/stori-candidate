package api

import (
	"os"
	"testing"

	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/g-ton/stori-candidate/mail"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store, mail mail.Mail) *Server {
	server := NewServer(store, mail)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
