package usecase

import (
	"RouletteGo/domain"
	"context"
)

type rouletteUsecase struct {
	rouletteRepo domain.RouletteRepository
}

// NewRouletteUsecase ...業務邏輯被新增出來時，實作介面
func NewRouletteUsecase(rouletteRepo domain.RouletteRepository) domain.RouletteUsecase {
	return &rouletteUsecase{
		rouletteRepo,
	}
}

func (r rouletteUsecase) GetByID(ctx context.Context, id string) (*domain.Bet, error) {
	panic("implement me")
}

func (r rouletteUsecase) Store(ctx context.Context, d *domain.Bet) error {
	panic("implement me")
}
