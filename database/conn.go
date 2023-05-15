package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"score-summary-backend/configs"
	"time"
)

var DB *gorm.DB

func init() {
	var err error

	MysqlS := configs.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlS.UserName, MysqlS.Password, MysqlS.IpHost, MysqlS.DbName)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// 设置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get db instance: %v", err))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 验证是否连接成功
	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping database: %v", err))
	}
}
