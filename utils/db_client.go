package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gitlab.hho-inc.com/devops/flowctl-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DBClient struct {
	db *gorm.DB
	sqlDB *sql.DB
}

func NewDBClient() *DBClient {
	dsn := "hhodb:hhodb@2022@tcp(mysql8.hho-inc.com)/oops?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	// 设置最大的空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大的打开的连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大可重用的时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &DBClient{
		db: db,
		sqlDB: sqlDB,
	}
}

func (db *DBClient) DBStats() {
	//stats := new(sql.DBStats)
	stats, _ := json.Marshal(db.sqlDB.Stats())
	fmt.Printf(string(stats))
}

func (db *DBClient) DBPing() {
	if err := db.sqlDB.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("Pong")
	}
}

func (db *DBClient) AutoMigrate() {
	db.db.AutoMigrate(&models.CDHistory{})
	db.db.AutoMigrate(&models.CDStatus{})
}

func (db *DBClient) DBInsertHistory(rowHistory *models.CDHistory) {
	db.db.Debug().Create(rowHistory)
}

func (db *DBClient) DBQueryHistory(rowHistory *models.CDHistory) {
	//result := db.db.Debug().Where("app = ? And env = ?", rowHistory.App, rowHistory.Env).Find(rowHistory)
	result := []map[string]interface{}{}
	db.db.Debug().Model(rowHistory).Where("app = ? AND env = ?", rowHistory.App, rowHistory.Env).Find(&result)
	d, _ := json.Marshal(result)
	fmt.Println(string(d))
	fmt.Println(len(result))
}

func (db *DBClient) DBInsertOrUpdateStatus(rowStatus *models.CDStatus) {
	result := []map[string]interface{}{}
	db.db.Debug().Model(rowStatus).Where("app = ? And env = ?", rowStatus.App, rowStatus.Env).Find(&result)
	if len(result) == 1 {
		db.db.Debug().Model(rowStatus).Where("app = ? And env = ?", rowStatus.App, rowStatus.Env).Updates(map[string]interface{}{
			"commit_id": rowStatus.CommitID,
			"branch": rowStatus.Branch,
			"image_tag": rowStatus.ImageTag,
			"image_url": rowStatus.ImageUrl,
			"git_url": rowStatus.GitUrl,
		})
	} else if len(result) == 0 {
		db.db.Debug().Create(rowStatus)
	} else {
		panic(fmt.Sprintf("%s 应用， %s 环境，查询结果超过一个", rowStatus.App, rowStatus.Env))
	}
}

//func main()  {
//	cli := NewDBClient()
//	t := new(models.CDHistory)
//	t.App = "aquaman-cart"
//	t.Env = "daily"
//	cli.DBQueryHistory(t)
//}