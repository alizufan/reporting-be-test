package merchant

import (
	"context"
	"errors"
	"reporting/schema"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/rotisserie/eris"
)

type MerchantRepository interface {
	Find(ctx context.Context, fil FindFilter) (*schema.Merchant, error)
}

type Merchant struct {
	DB rel.Repository
}

type FindFilter struct {
	ID           uint64
	UserID       uint64
	MerchantName string
}

func (f FindFilter) IsValidID() bool {
	return f.ID > 0
}

func (f FindFilter) IsValidUserID() bool {
	return f.UserID > 0
}

func (f FindFilter) IsValidMerchantName() bool {
	return len(f.MerchantName) > 0
}

func (u *Merchant) Find(ctx context.Context, fil FindFilter) (*schema.Merchant, error) {
	q := []rel.Querier{}

	if fil.IsValidID() {
		q = append(q, where.Eq("id", fil.ID))
	}

	if fil.IsValidUserID() {
		q = append(q, where.Eq("user_id", fil.UserID))
	}

	if fil.IsValidMerchantName() {
		q = append(q, where.Eq("merchant_name", fil.MerchantName))
	}

	res := schema.Merchant{}
	if err := u.DB.Find(ctx, &res, q...); err != nil {
		msg := "record not found"
		if !errors.Is(err, rel.ErrNotFound) {
			msg = "an error occurred"
		}
		return nil, eris.Wrapf(err, "find merchant, %s", msg)
	}

	return &res, nil
}
