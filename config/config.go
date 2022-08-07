package config

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

type mysqlDB struct {
	Host       string
	Port       string
	DbUser     string
	DbName     string
	DbPassword string
}

type app struct {
	Host string
	Port string
}

type grpc struct {
	GRPCAddress string
}

var (
	DB         *gorm.DB
	Mysql      mysqlDB
	App        app
	GRPC       grpc
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	SecretKey  string
)

func init() {
	iniPath := "config/config.ini"

	args := os.Args
	if len(args) > 1 {
		iniPath = args[1]
	}

	//loading
	iniFile, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalf("Load %v error %v \n", iniPath, err.Error())
		os.Exit(1)
	}

	//app config match ini file
	app := iniFile.Section("app")
	App.Host = app.Key("Host").String()
	App.Port = app.Key("Port").String()

	//grpc config
	grpc := iniFile.Section("grpc")
	GRPC.GRPCAddress = grpc.Key("GRPCAddress").String()
	// GRPC.Host = grpc.Key("Host").String()
	// GRPC.Port = grpc.Key("Port").String()

	//mysql
	mysql := iniFile.Section("mysql")
	Mysql.Host = mysql.Key("MysqlHost").String()
	Mysql.Port = mysql.Key("MysqlPort").String()
	Mysql.DbUser = mysql.Key("MysqlUser").String()
	Mysql.DbName = mysql.Key("MysqlDBName").String()
	Mysql.DbPassword = mysql.Key("MysqlPassword").String()

	//load rsa keys
	rsa := iniFile.Section("rsa")
	prvKey := rsa.Key("Private_Key_File").String()
	prvData, err := ioutil.ReadFile(prvKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	prv, err := jwt.ParseRSAPrivateKeyFromPEM(prvData)
	if err != nil {
		log.Fatal(err.Error())
	}
	PrivateKey = prv // assign privatekey pem

	pubKey := rsa.Key("Public_Key_File").String()
	pubData, err := ioutil.ReadFile(pubKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubData)
	if err != nil {
		log.Fatal(err.Error())
	}
	PublicKey = pub

	//Secret Key
	SecretKey = rsa.Key("Secret_Key").String()

}
