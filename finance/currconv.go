package finance

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/apierror"
	"github.com/renato0307/learning-go-lib/finance"
)

type getCurrConvOutput struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	Amount          float64 `json:"amount"`
	ConvertedAmount float64 `json:"converted_amount"`
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
			msg := "error: 'from' parameter is required"
			c.JSON(http.StatusBadRequest, apierror.New(msg))
			return
		}

		if to == "" {
			msg := "error: 'to' parameter is required"
			c.JSON(http.StatusBadRequest, apierror.New(msg))
			return
		}

		if amount == "" {
			msg := "error: 'amount' parameter is required"
			c.JSON(http.StatusBadRequest, apierror.New(msg))
			return
		}

		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			msg := "error: 'amount' is not a valid number"
			c.JSON(http.StatusBadRequest, apierror.New(msg))
			return
		}

		convertAmount, err := f.ConvertCurrency(from, to, amountFloat)
		if err != nil {
			msg := fmt.Sprintf("error converting the currency: %s", err.Error())
			c.JSON(http.StatusInternalServerError, apierror.New(msg))
			return
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
