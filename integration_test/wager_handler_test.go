//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"hoang/m/models"
	"net/http"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/require"
)

func placeWager() *models.Wager {
	wager := &models.Wager{
		TotalWagerValue:   100,
		Odds:              21,
		SellingPercentage: 30,
		SellingPrice:      35.213423,
	}
	models.CreateWager(wager)
	wager.Get()
	return wager
}

func TestPlaceWagerInvalidCases(t *testing.T) {
	body := []string{
		`{
			"total_wager_value": 0,
			"odds": 23,
			"selling_percentage": 4,
			"selling_price": 50
		}`,
		`{
			"total_wager_value": "abc",
			"odds": -1,
			"selling_percentage": -1,
			"selling_price": 100
		}`,
		`{
			"total_wager_value": 100,
			"odds": 4,
			"selling_percentage": -90,
			"selling_price": 5
		}`,
		`{
			"total_wager_value": 100,
			"odds": 20,
			"selling_percentage": 60,
			"selling_price": 30
		}`,
		`{
			"total_wager_value": "abc",
			"odds": "abc",
			"selling_percentage": "abc",
			"selling_price": "abc"
		}`,
	}

	for _, b := range body {
		resp, _, errs := gorequest.New().Post(testServer.URL + "/wagers").
			Send(b).
			End()
		require.Equal(t, 0, len(errs))
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPlaceWagerValidCase(t *testing.T) {
	body := `{
		"total_wager_value": 100,
		"odds": 20,
		"selling_percentage": 60,
		"selling_price": 80.3565645
	}`
	resp, bodyStr, errs := gorequest.New().Post(testServer.URL + "/wagers").
		Send(body).
		End()
	require.Equal(t, 0, len(errs))
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var bodyResp models.Wager
	json.Unmarshal([]byte(bodyStr), &bodyResp)
	require.Equal(t, int64(100), bodyResp.TotalWagerValue)
	require.Equal(t, int64(20), bodyResp.Odds)
	require.Equal(t, int64(60), bodyResp.SellingPercentage)
	require.Equal(t, float64(80.36), bodyResp.SellingPrice)
	require.Equal(t, float64(80.36), bodyResp.CurrentSellingPrice)
	fmt.Println(bodyStr)
}

func TestListWagersInValid(t *testing.T) {
	params := [][]string{
		{"abc", "1"},
		{"121", "xyz"},
		{"abc", "xyz"},
	}
	for _, p := range params {
		page := p[0]
		limit := p[1]
		resp, _, errs := gorequest.New().Get(fmt.Sprintf("%s/wagers?page=%s&limit=%s", testServer.URL, page, limit)).End()
		require.Equal(t, 0, len(errs))
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestListWagers(t *testing.T) {
	n := 20
	for i := 0; i < n; i++ {
		placeWager()
	}
	page := 2
	limit := 5
	resp, bodyStr, errs := gorequest.New().Get(fmt.Sprintf("%s/wagers?page=%d&limit=%d", testServer.URL, page, limit)).End()
	require.Equal(t, 0, len(errs))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	wagers := []*models.Wager{}
	json.Unmarshal([]byte(bodyStr), &wagers)
	require.Equal(t, 5, len(wagers))
}
