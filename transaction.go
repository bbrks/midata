package midata

import (
	"bytes"
	"encoding/csv"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
)

// transactionRegexp is used in isTransaction to determine if the line is a transaction
var transactionRegexp = regexp.MustCompile(`\d{2}/\d{2}/\d{4},[A-Z]+,.+(,(\-|\+)?\d+\.\d{2}){2}`)

// Transactions are held in a midata file
type Transactions []*Transaction

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Transactions) Less(i, j int) bool { return t[i].Date.Before(t[j].Date) }

// Transaction is a single entry in a midata file
type Transaction struct {
	Date     time.Time
	Type     string
	Merchant string

	Debit   decimal.Decimal
	Balance decimal.Decimal
}

// isTransaction returns true if the given bytes
// match the signature of a midata transaction.
func isTransaction(b []byte) bool {
	return transactionRegexp.Match(b)
}

func parseTransaction(b []byte) (*Transaction, error) {
	csv := csv.NewReader(bytes.NewReader(b))

	data, err := csv.Read()
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(dateLayout, data[0])
	if err != nil {
		return nil, err
	}

	debit, err := decimal.NewFromString(data[3])
	if err != nil {
		return nil, err
	}

	balance, err := decimal.NewFromString(data[4])
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Date:     date,
		Type:     data[1],
		Merchant: data[2],
		Debit:    debit,
		Balance:  balance,
	}, nil
}
