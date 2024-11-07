package seeders

import (
	"github.com/sirupsen/logrus"
	"github.com/snykk/transaction-api/internal/constants"
	"github.com/snykk/transaction-api/internal/datasources/records"
	"github.com/snykk/transaction-api/pkg/helpers"
	"github.com/snykk/transaction-api/pkg/logger"
)

var ProductData []records.Product

func init() {
	var err error
	pass, err = helpers.GenerateHash("12345")
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}

	ProductData = []records.Product{
		{
			Name:        "Keyboard",
			Description: "lorem ipsum dolor sit amet",
			Price:       240000,
			Stock:       23,
		},
		{
			Name:        "Mouse",
			Description: "lorem ipsum dolor sit amet",
			Price:       200000,
			Stock:       100,
		},
		{
			Name:        "Wifi Adapter",
			Description: "lorem ipsum dolor sit amet",
			Price:       150000,
			Stock:       234,
		},
		{
			Name:        "Computer Fan",
			Description: "lorem ipsum dolor sit amet",
			Price:       100000,
			Stock:       200,
		},
	}
}
