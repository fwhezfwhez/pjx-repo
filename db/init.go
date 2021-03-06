package db

import (
	"database/sql"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // 导入postgres
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *gorm.DB
var PQDB *sql.DB

const (
	host     = "localhost"
	user     = "postgres"
	dbname   = "test"
	sslmode  = "disable"
	password = "password"
	port     = "5432"
)

func init() {
	// 初始化数据库orm连接
	dbConfig := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s port=%s",
		host,
		user,
		dbname,
		sslmode,
		password,
		port,
	)
	log.Println(dbConfig)
	db, err := gorm.Open("postgres",
		dbConfig,
	)
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(10 * time.Second)
	db.DB().SetMaxIdleConns(30)
	if err != nil {
		panic(err)
	} else {
		DB = db
	}
	if e := DB.DB().Ping(); e != nil {
		panic(e)
	}

	// 自动重连，每60秒ping一次，失败时自动重连，重连间隔依次为3s,3s,15s,30s,60s,60s,60s.....
	go func(dbConfig string) {
		var intervals = []time.Duration{3 * time.Second, 3 * time.Second, 15 * time.Second, 30 * time.Second, 60 * time.Second,
		}
		for {
			time.Sleep(60 * time.Second)
			if e := DB.DB().Ping(); e != nil {
			L:
				for i := 0; i < len(intervals); i++ {
					e2 := RetryHandler(3, func() (bool, error) {
						var e error
						DB, e = gorm.Open("postgres", dbConfig)
						if e != nil {
							return false, errorx.Wrap(e)
						}
						return true, nil
					})
					if e2 != nil {
						fmt.Println(e.Error())
						time.Sleep(intervals[i])
						if i == len(intervals)-1 {
							i--
						}
						continue
					}
					DB.SingularTable(true)
					DB.LogMode(true)
					DB.DB().SetConnMaxLifetime(10 * time.Second)
					DB.DB().SetMaxIdleConns(30)
					break L
				}

			}
		}
	}(dbConfig)

	// 初始化pq驱动,用于CopyIn
	PQDB, err = sql.Open("postgres", dbConfig)
	PQDB.SetConnMaxLifetime(10 * time.Second)
	PQDB.SetMaxIdleConns(1)
	if err != nil {
		panic(err)
	}
}
