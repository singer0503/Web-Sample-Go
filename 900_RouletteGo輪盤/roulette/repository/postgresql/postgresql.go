package postgresql

import (
	"RouletteGo/domain"
	"context"
	"database/sql"
)

type postgresqlRouletteRepository struct {
	db *sql.DB
}

// NewPostgresqlRouletteRepository ... 實作其他後端介面
func NewPostgresqlRouletteRepository(db *sql.DB) domain.RouletteRepository {
	return &postgresqlRouletteRepository{db}
}

func (p postgresqlRouletteRepository) GetByID(ctx context.Context, id string) (*domain.Bet, error) {
	panic("implement me")
}

func (p postgresqlRouletteRepository) Store(ctx context.Context, d *domain.Bet) error {
	panic("implement me")
}
