package midata

import (
	"log"
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var m *Midata

func TestMarshal(t *testing.T) {
	assert.NotNil(t, m)
	assert.NotNil(t, m.Transactions)
	assert.NotNil(t, m.Overdraft)
}

func TestDaysInOverdraft(t *testing.T) {
	assert.Equal(t, m.DaysInOverdraft(), 0)
}

func TestTotalIncome(t *testing.T) {
	assert.Equal(t, 0, m.TotalIncome().Cmp(decimal.NewFromFloat(5.00)))
}

func TestTotalExpenses(t *testing.T) {
	assert.Equal(t, 0, m.TotalExpenses().Cmp(decimal.NewFromFloat(295)))
}

func TestTotalNet(t *testing.T) {
	assert.Equal(t, 0, m.TotalNet().Cmp(decimal.NewFromFloat(-290)))
}

func init() {
	f, err := os.Open("example/spec.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m, err = Marshal(f)
	if err != nil {
		log.Fatal(err)
	}
}
