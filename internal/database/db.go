package database

import (
	"chat/internal/config"
	"chat/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	config.LoadCfg()
	dsn := "host=" + config.AppConfig.PostgresHost +
		" user=" + config.AppConfig.PostgresUser +
		" password=" + config.AppConfig.PostgresPassword +
		" dbname=" + config.AppConfig.PostgresDB +
		" port=" + config.AppConfig.PostgresPort

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Println(err)
		return
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Println(err)
		return
	}
}
