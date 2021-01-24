package domain

import "context"

// 這是一個定義 業務邏輯層(Usecase) 以及 資料庫層（Repository）的介面
// 以功能別去區分，這個是輪盤專用的 interface

// Bet 物件 ...
type Bet struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string
}

// RouletteUsecase .. 業務邏輯介面
type RouletteUsecase interface {
	GetByID(ctx context.Context, id string) (*Bet, error)
	Store(ctx context.Context, d *Bet) error
}

// RouletteRepository ... 與其他後端做介接的介面
type RouletteRepository interface {
	GetByID(ctx context.Context, id string) (*Bet, error)
	Store(ctx context.Context, d *Bet) error
}
