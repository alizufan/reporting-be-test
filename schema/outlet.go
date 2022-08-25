package schema

import "time"

type Outlet struct {
	ID         uint64    `json:"id" db:"id"`
	MerchantID uint64    `json:"merchant_id" db:"merchant_id"`
	OutletName string    `json:"outlet_name" db:"outlet_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  uint64    `json:"created_by" db:"created_by"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy  uint64    `json:"updated_by" db:"updated_by"`
}

func (Outlet) Table() string {
	return "Outlets"
}
