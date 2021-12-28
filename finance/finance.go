package finance

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/finance"
)

type getCurrConvOutput struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	Amount          float64 `json:"amount"`
	ConvertedAmount float64 `json:"converted_amount"`
}

// SetRouterGroup defines all the routes for the finance functions
func SetRouterGroup(f finance.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	financeGroup := base.Group("/finance")
	{
		financeGroup.GET("/currconv", getCurrConv(f))
		// Add here more functions in the finance category
	}

	return financeGroup
}

// getCurrConv handles the currency conversion request.
//
// The request requires the from, to and amount parameters in the query string.
// It returns HTTP 200 on success.
// Returns HTTP 400 if there is a missing parameter.
// Returns HTTP 500 if there is another error.
func getCurrConv(f finance.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		from := c.Query("from")
		to := c.Query("to")
		amount := c.Query("amount")

		if from == "" {
			c.JSON(http.StatusBadRequest, "error: 'from' parameter is required")
			return
		}

		if to == "" {
			c.JSON(http.StatusBadRequest, "error: 'to' parameter is required")
			return
		}

		if amount == "" {
			c.JSON(http.StatusBadRequest, "error: 'amount' parameter is required")
			return
		}

		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "error: 'amount' is not a valid number")
			return
		}

		convertAmount, err := f.ConvertCurrency(from, to, amountFloat)
		if err != nil {
			err = fmt.Errorf("error converting the currency: %s", err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		output := getCurrConvOutput{
			From:            from,
			To:              to,
			Amount:          amountFloat,
			ConvertedAmount: convertAmount,
		}

		c.JSON(http.StatusOK, output)
	}
}
