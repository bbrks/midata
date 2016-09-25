package midata

import (
	"bufio"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

const (
	dateLayout = "02/01/2006"

	// The following are possible CSV headers seen in midata files.
	specHeader   = "Date,Type,Merchant/Description,Debit/Credit,Balance"
	lloydsHeader = "Transaction Date,Transaction Type,Merchant/Description,Debit/Credit,Balance"
)

var (
	// ErrInvalidMidata is returned when data
	// does not match the expected midata format.
	ErrInvalidMidata = errors.New("invalid midata format")

	headerRegexp = regexp.MustCompile(`(` + specHeader + "|" + lloydsHeader + `)`)
)

// Marshal takes an io.Reader and attempts
// to parse it into a midata struct.
func Marshal(r io.Reader) (*Midata, error) {
	var (
		line int
		m    Midata
	)
	s := bufio.NewScanner(r)

	for s.Scan() {
		line++
		b := s.Bytes()

		if isTransaction(b) {
			t, err := parseTransaction(b)
			if err != nil {
				return nil, err
			}

			m.Transactions = append(m.Transactions, t)
		} else if isOverdraft(b) {
			o, err := parseOverdraft(b)
			if err != nil {
				return nil, err
			}

			m.Overdraft = *o
		} else if line == 1 && isHeader(b) {
			continue
		} else {
			return nil, errors.Wrapf(ErrInvalidMidata, "line %d - %s", line, s.Text())
		}

	}

	return &m, nil
}

// isHeader returns true if the given bytes match
// the signature of a midata header.
func isHeader(b []byte) bool {
	return headerRegexp.Match(b)
}
