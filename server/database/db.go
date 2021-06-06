package database

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"github.com/wxmsummer/AirConditioner/server/scheduler"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

func InitDB() (*gorm.DB, error) {
	viper.SetConfigName("mysqlConfig")   //  设置配置文件名 (不带后缀)
	viper.AddConfigPath(".") 			// 比如添加当前目录
	err := viper.ReadInConfig()       		// 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf := map[string] string {
		"host":     viper.GetString("mysql.host"),
		"port":     viper.GetString("mysql.port"),
		"user":     viper.GetString("mysql.user"),
		"database": viper.GetString("mysql.database"),
		"password": viper.GetString("mysql.password"),
	}

	//db
	db, err := CreateConnection(conf)
	if err != nil {
		log.Fatalf("connection error : %v \n", err)
		return nil, err
	}

	// 初始化一些空调数据
	initData(db)
	scheduler.SchedulerDB = db

	return db, err
}

// db服务：db初始化：如果user表不存在则创建，然后连接数据库
func CreateConnection(conf map[string]string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(strings.Join(
		[]string {
			conf["user"], ":", conf["password"],						// user:pass
			"@tcp(", conf["host"], ":", conf["port"], ")/",			// @tcp(host:port)/
			conf["database"], "?charset=utf8mb4&parseTime=True&loc=Local",	// dbname?charset=utf8mb4&parseTime=True&loc=Local
		}, ""),
	), &gorm.Config {},
	)
}

// 初始化空调数据
func initData(db *gorm.DB) {
	air := model.AirConditioner{}
	airOrm := repository.AirConditionerOrm{Db: db}
	for i := 1001; i < 1006; i++ {
		air.RoomNum = i
		airOrm.Create(&air)
	}

	db.AutoMigrate(model.Admin{}, model.Fee{}, model.AirConditioner{}, model.Report{}, model.User{})


	fmt.Println("空调数据库记录初始化成功！")
}