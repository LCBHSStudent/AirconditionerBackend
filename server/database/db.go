package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"github.com/wxmsummer/AirConditioner/server/model"
	"log"
)

func InitDB() (*gorm.DB, error) {
	viper.SetConfigName("dbConfig")   //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("./database") // 比如添加当前目录
	err := viper.ReadInConfig()       // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf := map[string]interface{}{
		"host":     viper.Get("mysql.host"),
		"port":     viper.Get("mysql.port"),
		"user":     viper.Get("mysql.user"),
		"database": viper.Get("mysql.database"),
		"password": viper.Get("mysql.password"),
	}

	//db
	db, err := CreateConnection(conf)
	if err != nil {
		log.Fatalf("connection error : %v \n", err)
		return nil, err
	}
	db.Table("air_conditioners").CreateTable(&model.AirConditioner{})
	db.Table("fees").CreateTable(&model.Fee{})

	return db, err
}

// db服务：db初始化：如果user表不存在则创建，然后连接数据库
func CreateConnection(conf map[string]interface{}) (*gorm.DB, error) {
	host := conf["host"]
	port := conf["port"]
	user := conf["user"]
	dbName := conf["database"]
	password := conf["password"]
	return gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	),
	)
}
