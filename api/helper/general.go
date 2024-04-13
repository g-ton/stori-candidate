package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/g-ton/stori-candidate/mail"
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

		date := strings.Split(trans.Date, "/")
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

type TransactionByFile struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Transaction float64 `json:"transaction"`
}

func ProcessFile(filePath string) ([]db.Transaction, error) {
	//contentFile, err := os.ReadFile("./files/txns.csv")
	contentFile, err := os.ReadFile(filePath)
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

		// Converting ID from string to int64
		id, _ := strconv.ParseInt(row[0], 10, 64)
		transactions = append(transactions, db.Transaction{
			ID:          id,
			Date:        row[1],
			Transaction: transValue,
		})
	}

	return transactions, nil
}

func ProcessTemplateEmailForTransaction(tr TransactionResult, mails []string, mail mail.Mail) error {
	// Parsing the HTML template with the content for the email
	//absRootPath := util.GetAbsRootPath()
	//path := filepath.Join(absRootPath, "files", "stori-template.html")
	path := ""
	if os.Getenv("FOO") == "1" {
		path = "../files/stori-template.html"
	} else if os.Getenv("FOO") == "2" {
		path = "../../files/stori-template.html"
	} else {
		path = "files/stori-template.html"
	}
	t, err := template.ParseFiles(path)
	if err != nil {
		fmt.Println("error parsing the HTML template", err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Stori Challenge by Damián \n%s\n\n", mimeHeaders)))

	monthsSlice := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	transactionsPerMonth := []string{}
	for _, month := range monthsSlice {
		if v, ok := tr.Months[month]; ok {
			transactionsPerMonth = append(transactionsPerMonth, fmt.Sprintf("Number of transactions in %v: %v\n", month, v))
		}
	}

	t.Execute(&body, struct {
		TotalBalance     float64
		AvgCredit        float64
		AvgDebit         float64
		NumTransPerMonth []string
	}{
		TotalBalance:     tr.Data["totalBalance"],
		AvgCredit:        tr.Data["avgCredit"],
		AvgDebit:         tr.Data["avgDebit"],
		NumTransPerMonth: transactionsPerMonth,
	})

	err = mail.SendMail(mails, body.Bytes())
	if err != nil {
		fmt.Println("error sending email", err)
		return err
	}

	return nil
}
