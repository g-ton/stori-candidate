package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/g-ton/stori-candidate/api"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/g-ton/stori-candidate/mail"
	"github.com/g-ton/stori-candidate/util"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
)

func main() {
	os.Setenv("FOO", "0")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	mail := mail.NewMail(config)
	server := api.NewServer(store, mail)

	// Enable CORS (only for test purpouses)
	// ref: https://github.com/gin-contrib/cors
	server.Router.Use(cors.Default())

	// Routes (actions) for account
	server.Router.POST("/accounts", server.CreateAccount)
	server.Router.GET("/accounts/:id", server.GetAccount)
	server.Router.GET("/accounts", server.ListAccounts)
	// Routes (actions) for transactions
	server.Router.POST("/transactions", server.CreateTransaction)
	server.Router.GET("/transactions/:id", server.GetTransaction)
	server.Router.GET("/transactions", server.ListTransactions)
	server.Router.POST("/sendSummaryInfoByDB", server.GetSummaryInfoByDB)
	server.Router.POST("/sendSummaryInfoByFile", server.GetSummaryInfoByFile)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
}
