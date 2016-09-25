package midata

import (
	"bytes"
	"encoding/csv"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
)

// overdraftRegexp is used in isTransaction to determine if the line is a transaction
var overdraftRegexp = regexp.MustCompile(`Arranged (O|o)verdraft (L|l)imit,\d{2}/\d{2}/\d{4},\d+`)

// Overdraft information is appended to the midata file
type Overdraft struct {
	Date   time.Time
	Amount decimal.Decimal
}

// isOverdraft returns true if the given bytes
// match the signature of a midata overdraft section.
func isOverdraft(b []byte) bool {
	return overdraftRegexp.Match(b)
}

func parseOverdraft(b []byte) (*Overdraft, error) {
	csv := csv.NewReader(bytes.NewReader(b))

	data, err := csv.Read()
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(dateLayout, data[1])
	if err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(data[2])
	if err != nil {
		return nil, err
	}

	return &Overdraft{
		Date:   date,
		Amount: amount,
	}, nil
}
