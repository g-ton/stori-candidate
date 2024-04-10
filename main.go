package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/g-ton/stori-candidate/internal/util"
)

func main_() {
	//fmt.Println("Say hello")

	in := `first_name,last_name,username
	"Rob","Pike","rob"
	"Ken","Thompson","ken"
	"Robert","Griesemer","gri"
	`
	r := csv.NewReader(strings.NewReader(in))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)

	// _env := helpers.GetEnvDB()

	// // Run database (Initialazing) with env config
	// // (Database Service)
	// dbService := database.New(_env)
	// // (Web Server Service)
	// webService := websrv.New()

	// // Initialazing the handler to controll the web actions (Passing database and web server services)
	// wedHandler := webHdl.NewHTTPHandlerGin(dbService, webService)
	// router := wedHandler.Router
	// // Enable CORS (only for test purpouses)
	// // ref: https://github.com/gin-contrib/cors
	// router.Use(cors.Default())

	// // Routes (actions)
	// router.POST("/sendInfo", wedHandler.SendSummaryInformation)
	// wedHandler.Run()
}

func main() {
	transactions, _ := util.ProcessFile()
	r := util.GetSummaryInfo(transactions)

	fmt.Println("totalBalance: ", r.Data["totalBalance"])
	fmt.Println("avgCredit: ", r.Data["avgCredit"])
	fmt.Println("avgDebit: ", r.Data["avgDebit"])

	monthsSlice := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	for _, month := range monthsSlice {
		if v, ok := r.Months[month]; ok {
			fmt.Printf("Number of transactions in %v: %v\n", month, v)
		}
	}

}
