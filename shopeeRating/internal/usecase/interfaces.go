// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/TranQuocToan1996/shopeerating/internal/entity"
)

type (
	Rating interface {
		GetRatings(context.Context, string) (*entity.RatingResp, error)
	}

	ShopeeAPI interface {
		GetRatings(context.Context, string) (*entity.RatingResp, error)
	}
)
