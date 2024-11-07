package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/snykk/transaction-api/cmd/seed/seeders"
	"github.com/snykk/transaction-api/internal/config"
	"github.com/snykk/transaction-api/internal/constants"
	"github.com/snykk/transaction-api/internal/utils"
	"github.com/snykk/transaction-api/pkg/logger"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

func main() {
	db, err := utils.SetupPostgresConnection()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}
	defer db.Close()

	logger.Info("seeding...", logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})

	seeder := seeders.NewSeeder(db)
	err = seeder.UserSeeder(seeders.UserData)
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}

	logger.Info("seeding success!", logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
}
