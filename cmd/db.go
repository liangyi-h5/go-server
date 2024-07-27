package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/log"
)

// var db *sql.DB // 声明全局 db 变量
func initDb() *sql.DB {
	// 只读
	dsn := os.Getenv("READDB")
	// dsn := "username:password@tcp(127.0.0.1:3306)/yourdatabase?charset=utf8mb4&parseTime=True&loc=Local"

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(25)                 // 最大闲置连接数
	db.SetConnMaxLifetime(time.Minute * 5) // 连接最大生命周期
	// 检查连接是否有效
	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		log.Error(err)
	}

	log.Info("Connected to the database successfully!")
	return db
}

// var db *sql.DB // 声明全局 db 变量
func createWriteDB() *sql.DB {

	// 只读
	dsn := os.Getenv("WRITEDB")

	// 打开数据库连接
	WriteDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// 检查连接是否有效
	err = WriteDB.Ping()
	if err != nil {
		// log.Fatal(err)
		log.Error(err)
	}

	log.Info("Connected to the database successfully!")
	return WriteDB
}
