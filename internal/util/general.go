package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/g-ton/stori-candidate/db"
)

type TransactionResult struct {
	Months map[string]int
	Data   map[string]float64
}

func GetSummaryInfo(transactions []db.Transaction) TransactionResult {
	// Get the info from file or db and get the summary information
	if len(transactions) == 0 {
		return TransactionResult{}
	}

	// Creating the map to control operations
	numTransPerMonth := make(map[string]int)
	calcs := make(map[string]float64)
	calcs["totalBalance"] = 0.0
	calcs["avgDebit"] = 0.0
	calcs["avgCredit"] = 0.0
	iDebit := 0
	iCredit := 0
	totalDebit := 0.0
	totalCredit := 0.0

	for _, trans := range transactions {
		if trans.Date == "" {
			continue
		}

		date := strings.Trim(trans.Date, "/")
		month := date[0]

		switch string(month) {
		case "1":
			// January
			numTransPerMonth["January"] = numTransPerMonth["January"] + 1
		case "2":
			// February
			numTransPerMonth["February"] = numTransPerMonth["February"] + 1
		case "3":
			// March
			numTransPerMonth["March"] = numTransPerMonth["March"] + 1
		case "4":
			// April
			numTransPerMonth["April"] = numTransPerMonth["April"] + 1
		case "5":
			// May
			numTransPerMonth["May"] = numTransPerMonth["May"] + 1
		case "6":
			// June
			numTransPerMonth["June"] = numTransPerMonth["June"] + 1
		case "7":
			// July
			numTransPerMonth["July"] = numTransPerMonth["July"] + 1
		case "8":
			// August
			numTransPerMonth["August"] = numTransPerMonth["August"] + 1
		case "9":
			// September
			numTransPerMonth["September"] = numTransPerMonth["September"] + 1
		case "10":
			// October
			numTransPerMonth["October"] = numTransPerMonth["October"] + 1
		case "11":
			// November
			numTransPerMonth["November"] = numTransPerMonth["November"] + 1
		case "12":
			// December
			numTransPerMonth["December"] = numTransPerMonth["December"] + 1

		}

		calcs["totalBalance"] = calcs["totalBalance"] + trans.Transaction
		if trans.Transaction >= 0.0 {
			iCredit++
			// For credit transactions
			totalCredit += trans.Transaction
			calcs["avgCredit"] = totalCredit / float64(iCredit)
		} else {
			iDebit++
			// For debit transactions
			totalDebit += trans.Transaction
			calcs["avgDebit"] = totalDebit / float64(iDebit)
		}
	}

	return TransactionResult{
		Months: numTransPerMonth,
		Data:   calcs,
	}
}

func ProcessFile() ([]db.Transaction, error) {
	contentFile, err := os.ReadFile("./files/txns.csv")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	r := csv.NewReader(strings.NewReader(string(contentFile)))

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	transactions := []db.Transaction{}
	for _, row := range records {
		// Converting string from Transaction column into float64
		transValue, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}
		transactions = append(transactions, db.Transaction{
			ID:          row[0],
			Date:        row[1],
			Transaction: transValue,
		})
	}

	fmt.Println(records)
	return transactions, nil
}
