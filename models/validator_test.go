package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWagerValidator(t *testing.T) {
	wagers := []*Wager{
		{
			TotalWagerValue:   0,
			Odds:              23,
			SellingPercentage: 4,
			SellingPrice:      50,
		},
		{
			TotalWagerValue:   100,
			Odds:              -23,
			SellingPercentage: 4,
			SellingPrice:      50,
		},
		{
			TotalWagerValue:   100,
			Odds:              3,
			SellingPercentage: -70,
			SellingPrice:      78,
		},
		{
			TotalWagerValue:   100,
			Odds:              3,
			SellingPercentage: 70,
			SellingPrice:      57,
		},
		{
			TotalWagerValue:   100,
			Odds:              3,
			SellingPercentage: 70,
			SellingPrice:      80,
		},
	}
	cnt := 0
	for _, w := range wagers {
		wagerValidator := &WagerValidator{
			Wager: w,
		}
		if Validate(wagerValidator) {
			cnt++
		}
	}
	require.Equal(t, 1, cnt)
}

func TestPurchaseValidator(t *testing.T) {
	purchases := []*Purchase{
		{
			BuyingPrice: -1,
		},
		{
			BuyingPrice: 100,
		},
	}
	cnt := 0
	for _, p := range purchases {
		purchaseValidator := &PurchaseValidator{
			Purchase: p,
		}
		if Validate(purchaseValidator) {
			cnt++
		}
	}
	require.Equal(t, 1, cnt)
}
