package schema

import "time"

type Merchant struct {
	ID           uint64    `json:"id" db:"id"`
	UserID       uint64    `json:"user_id" db:"user_id"`
	MerchantName string    `json:"merchant_name" db:"merchant_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	CreatedBy    uint64    `json:"created_by" db:"created_by"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy    uint64    `json:"updated_by" db:"updated_by"`
}

func (Merchant) Table() string {
	return "Merchants"
}
