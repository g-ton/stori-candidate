package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/g-ton/stori-candidate/api"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/g-ton/stori-candidate/mail"
	"github.com/g-ton/stori-candidate/util"
	_ "github.com/lib/pq"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	os.Setenv("EMAIL_TEMPLATE_LEVEL", "0")
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

	// Routes (actions) for transactions
	server.Router.POST("/transactions", server.CreateTransaction)
	server.Router.GET("/transactions/:id", server.GetTransaction)
	server.Router.GET("/transactions", server.ListTransactions)

	ginLambda = ginadapter.New(server.Router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}