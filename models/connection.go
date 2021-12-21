package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
)

var Db *gorm.DB

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	user := utils.Viper.GetString("DB_USER")
	password := utils.Viper.GetString("DB_PASS")
	server := utils.Viper.GetString("DB_SERVER")
	database := "iglesias_develop"
	if utils.Viper.GetBool("PROD") {
		database = utils.Viper.GetString("DB_NAME")
	}
	port := utils.Viper.GetString("DB_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, server, port, database)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().In(tiempo.Local)
		},
	})
	if err != nil {
		utils.Log.Fatal("Error al conectar base de datos", err)
	}
	utils.Log.Info("Conectado a: " + dsn)
	if !utils.Viper.GetBool("PROD") {
		//migrarTablas()
	}
}

func migrarTablas() {
	//Poner tablas para migrar
	err := Db.AutoMigrate(&Usuario{}, &AdminMaster{}, &AdminParroquia{}, &Fiel{})
	if err != nil {
		utils.Log.Warn(err)
		utils.Log.Fatal("Error al migrar modelos")

	}
}

type Tabler interface {
	TableName() string
}
