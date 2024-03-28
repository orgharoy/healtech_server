package database

import (
	"log"

	model "github.com/orgharoy/healtech/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	//dsn := "root:Ed2H4be54A41hfCBcg346AgBBDa1bCDa@tcp(viaduct.proxy.rlwy.net:17031)/railway" //->Orgha
	//dsn := "root:CCE5EC-f6Abef2eCga1d1gAD5e3FBD12@tcp(viaduct.proxy.rlwy.net:48576)/railway?parseTime=true" // -> Saddam
	dsn := "root:admin1@3@tcp(localhost:3306)/healtech?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")

	DB.AutoMigrate(
		&model.User{},
		&model.Report{},
		&model.ReportGroup{},
		&model.Patient{},
		&model.PatientBill{},
		&model.BillDetails{},
	)

	log.Println("🚀 Connected Successfully to the Database")

	return nil
}
