package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

var ErrWagerNotFound = errors.New("no wager")
var ErrBuyingPriceGreaterThan = errors.New("buying_price is greater than current_selling_price")

type Purchase struct {
	bun.BaseModel `bun:"table:purchases,alias:p"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	WagerID     int64     `bun:"wager_id,notnull" json:"wager_id"`
	BuyingPrice float64   `bun:"buying_price,notnull" json:"buying_price"`
	BoughtAt    time.Time `bun:"bought_at,nullzero,notnull,default:current_timestamp" json:"bought_at"`
}

func CreatePurchase(p *Purchase) error {
	fail := func(err error) error {
		return fmt.Errorf("CreateOrder: %v", err)
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	wager := &Wager{ID: p.WagerID}
	err = wager.GetInTransaction(tx, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrWagerNotFound
		}
		return fail(err)
	}
	if wager.CurrentSellingPrice < p.BuyingPrice {
		return ErrBuyingPriceGreaterThan
	}
	err = wager.UpdateInTransaction(p, tx, ctx)
	if err != nil {
		return fail(err)
	}
	// Create Purchase
	err = tx.QueryRow(`INSERT INTO "purchases" (
		"wager_id",
		"buying_price"
	)
	VALUES (?, ?) RETURNING id, bought_at`, p.WagerID, p.BuyingPrice).Scan(&p.ID, &p.BoughtAt)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return err
}
