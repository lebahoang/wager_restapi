package handlers

import (
	"hoang/m/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func placeWager(c *gin.Context) {
	var wager models.Wager
	err := c.ShouldBindJSON(&wager)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	validator := &models.WagerValidator{
		Wager: &wager,
	}
	if !models.Validate(validator) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	w := &wager
	err = models.CreateWager(w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = w.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func buyWager(c *gin.Context) {
	wagerIDStr := c.Param("wager_id")
	wagerID, err := strconv.ParseInt(wagerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query param wager_id is invalid",
		})
		return
	}
	p := &models.Purchase{}
	err = c.ShouldBindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	p.WagerID = wagerID
	purchaseValidator := &models.PurchaseValidator{
		Purchase: p,
	}
	if !models.Validate(purchaseValidator) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is invalid",
		})
		return
	}
	err = models.CreatePurchase(p)
	if err != nil {
		if err == models.ErrWagerNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if err == models.ErrBuyingPriceGreaterThan {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusCreated, p)
}

func listWagers(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query param page is invalid",
		})
		return
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query param limit is invalid",
		})
		return
	}
	wagers, err := models.ListWager(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, wagers)
}
