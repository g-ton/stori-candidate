package api

import (
	"os"
	"testing"

	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	server := NewServer(store)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
