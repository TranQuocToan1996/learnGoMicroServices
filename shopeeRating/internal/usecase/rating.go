package usecase

import (
	"context"

	"github.com/TranQuocToan1996/shopeerating/config"
	"github.com/TranQuocToan1996/shopeerating/internal/entity"
	"github.com/TranQuocToan1996/shopeerating/internal/usecase/shopeeAPI/shopeev2"
)

type RatingUseCase struct {
	api ShopeeAPI
}

func New(cfg config.Config) *RatingUseCase {
	return &RatingUseCase{
		api: shopeev2.New(cfg),
	}
}

func (uc *RatingUseCase) GetRatings(ctx context.Context, rawURL string) (*entity.RatingResp, error) {
	return uc.api.GetRatings(ctx, rawURL)
}
