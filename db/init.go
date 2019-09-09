package db

import (
	"database/sql"
	"errorX"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // 导入postgres
	_ "github.com/lib/pq"
	"log"
	"shangraomajiang/config"
	"shangraomajiang/util/common/util"
	"time"
)

var DB *gorm.DB
var PQDB *sql.DB

func init() {
	// 初始化数据库orm连接
	c := config.GetConfig()

	dbConfig := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		c.GetString("db.host"),
		c.GetString("db.user"),
		c.GetString("db.dbname"),
		c.GetString("db.sslmode"),
		c.GetString("db.password"),
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

	go func(dbConfig string) {
		var intervals = []time.Duration{3 * time.Second, 3 * time.Second, 15 * time.Second, 30 * time.Second, 60 * time.Second,
			120 * time.Second, 360 * time.Second, 10 * time.Minute, 1 * time.Hour, 24 * time.Hour,
		}
		for {
			time.Sleep(30 * time.Second)
			if e := DB.DB().Ping(); e != nil {
			L:
				for i := 0; i < len(intervals); i++ {
					e2 := util.RetryHandler(3, func() (bool, error) {
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
