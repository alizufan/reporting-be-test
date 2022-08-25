package user

import (
	"context"
	"errors"
	"reporting/schema"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/rotisserie/eris"
)

type UserRepository interface{
	Find(ctx context.Context, fil FindFilter) (*schema.User, error) 
}


type User struct {
	DB rel.Repository
}

type FindFilter struct {
	ID       uint64
	Username string
}

func (f FindFilter) IsValidID() bool {
	return f.ID > 0
}

func (f FindFilter) IsValidUsername() bool {
	return len(f.Username) > 0
}

func (u *User) Find(ctx context.Context, fil FindFilter) (*schema.User, error) {
	q := []rel.Querier{}

	if fil.IsValidID() {
		q = append(q, where.Eq("id", fil.ID))
	}

	if fil.IsValidUsername() {
		q = append(q, where.Eq("user_name", fil.Username))
	}

	res := schema.User{}
	if err := u.DB.Find(ctx, &res, q...); err != nil {
		msg := "record not found"
		if !errors.Is(err, rel.ErrNotFound) {
			msg = "an error occurred"
		}
		return nil, eris.Wrapf(err, "find user, %s", msg)
	}

	return &res, nil
}
