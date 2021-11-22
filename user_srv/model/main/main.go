package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hash"
	"log"
	"mxshop_srvs/user_srv/model"
	"os"
	"time"
)

func main() {
	dsn := "root:142212@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdin, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println("数据库迁移失败：", err)
		return
	}
	type Options struct {
		SaltLen      int
		Iterations   int
		KeyLen       int
		HashFunction func() hash.Hash
	}

	// Using custom options
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	var users []model.User
	for x := 0; x < 10; x++ {
		user := model.User{
			NickName: fmt.Sprintf("hq%d", x),
			Mobile:   fmt.Sprintf("18281222%d", x),
			Password: newPassword,
		}
		//db.Save(&user)
		users = append(users, user)
	}
	result := db.CreateInBatches(users, 10)
	if result.RowsAffected == 0 {
		fmt.Println(result.Error)
	}

}
