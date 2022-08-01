package ds

import (
	"fmt"
	"miniproject/config"
	"miniproject/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSource struct {
	Mysql *gorm.DB
}

var DB *gorm.DB

func NewDataSource() *DataSource {

	host := config.Mysql.Host
	port := config.Mysql.Port
	name := config.Mysql.DbName
	password := config.Mysql.DbPassword
	user := config.Mysql.DbUser

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error on connecting to database!")
	}
	DB = db
	db.AutoMigrate(
		&model.User{},
		&model.VerifyCode{},
		&model.UserToken{},
		&model.MoneyEntries{},
		&model.MoneyTransfer{},
		&model.UserMoney{},
		&model.UserBankCard{},
		&model.UserCityTotalBankCard{},
	)

	return &DataSource{
		Mysql: db,
	}

}
