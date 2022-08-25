package schema

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID         uint64    `json:"id" db:"id"`
	MerchantID uint64    `json:"merchant_id" db:"merchant_id"`
	OutletID   uint64    `json:"outlet_id" db:"outlet_id"`
	BillTotal  float64   `json:"bill_total" db:"bill_total"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  uint64    `json:"created_by" db:"created_by"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy  uint64    `json:"updated_by" db:"updated_by"`
}

func (Transaction) Table() string {
	return "Transactions"
}

type CountTransactionReport struct {
	Count int `json:"count" db:"count"`
}

type TransactionReport struct {
	Date  time.Time `json:"date" db:"date"`
	Omzet string    `json:"omzet" db:"omzet"`
}

func (d *TransactionReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"date":  d.Date.Format("2006-01-02"),
		"omzet": d.Omzet,
	})
}
