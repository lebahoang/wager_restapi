//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"hoang/m/models"
	"net/http"
	"strconv"
	"sync"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/require"
)

func TestBuyWagerInvalidCases(t *testing.T) {
	body := []string{
		`{
			"buying_price": "dsad"
		}`,
		`{
			"buying_price": -223
		}`,
	}

	for _, b := range body {
		resp, _, errs := gorequest.New().Post(testServer.URL + "/buy/1").
			Send(b).
			End()
		require.Equal(t, 0, len(errs))
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestBuyWagerInvalidWagerID(t *testing.T) {
	body := `{"buying_price": 223}`
	for _, b := range body {
		resp, _, errs := gorequest.New().Post(testServer.URL + "/buy/dasd").
			Send(b).
			End()
		require.Equal(t, 0, len(errs))
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestBuyingPriceGreaterThan(t *testing.T) {
	wager := placeWager()
	body := `{"buying_price": 3000}`
	resp, _, errs := gorequest.New().Post(testServer.URL + "/buy/" + strconv.Itoa(int(wager.ID))).
		Send(body).
		End()
	require.Equal(t, 0, len(errs))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func TestBuyWagerValid(t *testing.T) {
	wager := placeWager()
	buyingPrice := 27
	body := fmt.Sprintf(`{"buying_price": %d}`, buyingPrice)
	resp, bodyStr, errs := gorequest.New().Post(testServer.URL + "/buy/" + strconv.Itoa(int(wager.ID))).
		Send(body).
		End()
	require.Equal(t, 0, len(errs))
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var bodyResp models.Purchase
	json.Unmarshal([]byte(bodyStr), &bodyResp)
	require.Equal(t, float64(buyingPrice), bodyResp.BuyingPrice)
	require.Equal(t, wager.ID, bodyResp.WagerID)
	require.NotNil(t, bodyResp.BoughtAt)
}

func TestMultipleBuyWagerValidSomeFailed(t *testing.T) {
	wager := placeWager()
	var wg sync.WaitGroup
	numberOfRequest := 50
	buyingPrice := 6
	expectedSuccessRequest := (int(wager.CurrentSellingPrice) - (int(wager.CurrentSellingPrice) % buyingPrice)) / buyingPrice
	if expectedSuccessRequest > numberOfRequest {
		expectedSuccessRequest = numberOfRequest
	}
	expectedFailRequest := numberOfRequest - expectedSuccessRequest
	rs := make([][]interface{}, numberOfRequest)
	for i := 0; i < numberOfRequest; i++ {
		wg.Add(1)
		go func(ind int) {
			defer wg.Done()
			body := fmt.Sprintf(`{"buying_price": %d}`, buyingPrice)
			resp, bodyStr, errs := gorequest.New().Post(testServer.URL + "/buy/" + strconv.Itoa(int(wager.ID))).
				Send(body).
				End()
			require.Equal(t, 0, len(errs))
			rs[ind] = []interface{}{resp, bodyStr}
		}(i)
	}
	wg.Wait()
	succeededRequests := 0
	failedRequests := 0
	for i := 0; i < numberOfRequest; i++ {
		resp := rs[i][0].(gorequest.Response)
		if resp.StatusCode == http.StatusCreated {
			succeededRequests++
		} else if resp.StatusCode == http.StatusBadRequest {
			failedRequests++
		}
	}
	require.Equal(t, expectedSuccessRequest, succeededRequests)
	require.Equal(t, expectedFailRequest, failedRequests)
}

func TestMultipleBuyWagerValidAllSucceeded(t *testing.T) {
	wager := placeWager()
	var wg sync.WaitGroup
	numberOfRequest := 20
	buyingPrice := 1
	expectedSuccessRequest := (int(wager.CurrentSellingPrice) - (int(wager.CurrentSellingPrice) % buyingPrice)) / buyingPrice
	if expectedSuccessRequest > numberOfRequest {
		expectedSuccessRequest = numberOfRequest
	}
	expectedFailRequest := numberOfRequest - expectedSuccessRequest
	rs := make([][]interface{}, numberOfRequest)
	for i := 0; i < numberOfRequest; i++ {
		wg.Add(1)
		go func(ind int) {
			defer wg.Done()
			body := fmt.Sprintf(`{"buying_price": %d}`, buyingPrice)
			resp, bodyStr, errs := gorequest.New().Post(testServer.URL + "/buy/" + strconv.Itoa(int(wager.ID))).
				Send(body).
				End()
			require.Equal(t, 0, len(errs))
			rs[ind] = []interface{}{resp, bodyStr}
		}(i)
	}
	wg.Wait()
	succeededRequests := 0
	failedRequests := 0
	for i := 0; i < numberOfRequest; i++ {
		resp := rs[i][0].(gorequest.Response)
		if resp.StatusCode == http.StatusCreated {
			succeededRequests++
		} else if resp.StatusCode == http.StatusBadRequest {
			failedRequests++
		}
	}
	require.Equal(t, expectedSuccessRequest, succeededRequests)
	require.Equal(t, expectedFailRequest, failedRequests)
}
