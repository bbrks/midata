package midata

import (
	"time"

	"github.com/shopspring/decimal"
)

// Midata holds midata information
type Midata struct {
	Transactions Transactions
	Overdraft    Overdraft
}

// DaysInOverdraft will return the total number of days with a negative balance.
func (md *Midata) DaysInOverdraft() int {
	m := make(map[time.Time]struct{})
	for _, t := range md.Transactions {
		if t.Balance.Cmp(decimal.Zero) < 0 {
			m[t.Date] = struct{}{}
		}
	}
	return len(m)
}

// TotalIncome will return the sum of credit for the account.
func (md *Midata) TotalIncome() decimal.Decimal {
	total := decimal.Zero
	for _, v := range md.Transactions {
		if v.Debit.Cmp(decimal.Zero) > 0 {
			total = total.Add(v.Debit)
		}
	}
	return total
}

// TotalExpenses will return the total amount of debit for the account.
func (md *Midata) TotalExpenses() decimal.Decimal {
	total := decimal.Zero
	for _, v := range md.Transactions {
		if v.Debit.Cmp(decimal.Zero) < 0 {
			total = total.Add(v.Debit.Abs())
		}
	}
	return total
}

// TotalNet will return the total net amount for the account.
func (md *Midata) TotalNet() decimal.Decimal {
	total := decimal.Zero
	for _, v := range md.Transactions {
		total = total.Add(v.Debit)
	}
	return total
}
