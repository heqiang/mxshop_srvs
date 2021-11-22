package initlialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"os"
	"time"
)

func Initdb() {
	// mxshop_user_srv
	conf := global.ServerConfig.MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password,
		conf.Host, conf.Port,
		conf.Db)
	newLogger := logger.New(
		log.New(os.Stdin, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("数据库打开失败：", err)
		return
	}
	err = global.DB.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println("数据库迁移失败：", err)
		return
	}
}
