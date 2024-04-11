package main

import (
	"database/sql"
	"log"

	"github.com/g-ton/stori-candidate/api"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/g-ton/stori-candidate/util"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

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
	server.Router.POST("/sendSummaryInfoByDB/:account_id", server.GetSummaryInfoByDB)
	server.Router.POST("/sendSummaryInfoByFile", server.GetSummaryInfoByFile)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
}

// func main_() {
// 	transactions, _ := util.ProcessFile()
// 	r := util.GetSummaryInfo(transactions)

// 	fmt.Println("totalBalance: ", r.Data["totalBalance"])
// 	fmt.Println("avgCredit: ", r.Data["avgCredit"])
// 	fmt.Println("avgDebit: ", r.Data["avgDebit"])

// 	monthsSlice := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
// 	for _, month := range monthsSlice {
// 		if v, ok := r.Months[month]; ok {
// 			fmt.Printf("Number of transactions in %v: %v\n", month, v)
// 		}
// 	}

// }
