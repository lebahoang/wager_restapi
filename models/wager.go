package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type Wager struct {
	bun.BaseModel `bun:"table:wagers,alias:w"`

	ID                  int64     `bun:"id,pk,autoincrement" json:"id"`
	TotalWagerValue     int64     `bun:"total_wager_value,notnull" json:"total_wager_value"`
	Odds                int64     `bun:"odds,notnull" json:"odds"`
	SellingPercentage   int64     `bun:"selling_percentage,notnull" json:"selling_percentage"`
	SellingPrice        float64   `bun:"selling_price,notnull" json:"selling_price"`
	CurrentSellingPrice float64   `bun:"current_selling_price,notnull" json:"current_selling_price"`
	PercentageSold      int64     `bun:"percentage_sold" json:"percentage_sold"`
	AmountSold          float64   `bun:"amount_sold" json:"amount_sold"`
	PlacedAt            time.Time `bun:"placed_at,nullzero,notnull,default:current_timestamp" json:"placed_at"`
}

func CreateWager(wager *Wager) error {
	err := db.QueryRow(`INSERT INTO "wagers" (
			"total_wager_value", "odds",
			"selling_percentage", "selling_price", "current_selling_price"
		)
		VALUES (?, ?, ?, ?, ?) RETURNING id`,
		wager.TotalWagerValue,
		wager.Odds,
		wager.SellingPercentage,
		wager.SellingPrice,
		wager.SellingPrice).Scan(&wager.ID)
	if err != nil {
		return err
	}
	return err
}

func ListWager(page int64, limit int64) ([]Wager, error) {
	wagers := []Wager{}
	offset := int64(0)
	if page > 1 {
		offset = (page - 1) * limit
	}
	ctx := context.Background()
	err := db.NewSelect().Model(&wagers).Limit(int(limit)).Offset(int(offset)).Scan(ctx)
	return wagers, err
}

func (w *Wager) Get() (int, error) {
	ctx := context.Background()
	cnt, err := db.NewSelect().Model(w).Where("id = ?", w.ID).ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return cnt, nil
	}
	return cnt, err
}

func (w *Wager) GetInTransaction(tx bun.Tx, ctx context.Context) error {
	return tx.NewSelect().Model(w).Where("id = ?", w.ID).For("UPDATE").Scan(ctx)
}

func (w *Wager) UpdateInTransaction(p *Purchase, tx bun.Tx, ctx context.Context) error {
	w.CurrentSellingPrice = w.CurrentSellingPrice - p.BuyingPrice
	w.AmountSold = w.AmountSold + p.BuyingPrice
	w.PercentageSold = w.PercentageSold + int64(((float64(100) * p.BuyingPrice) / float64(w.TotalWagerValue)))
	_, err := tx.NewUpdate().Model(w).WherePK().Exec(ctx)
	return err
}
