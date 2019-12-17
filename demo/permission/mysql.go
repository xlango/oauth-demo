package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func GetMysqlDb() (db *gorm.DB) {
	var err error
	//连接串
	db, err = gorm.Open("mysql", "root:123456@tcp(192.168.10.234:3306)/oauthdemo?charset=utf8")
	//defer db.Close()
	if err != nil {
		panic(err)
	}
	//设置最大空闲连接数和最大连接
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(100)
	//true:不使用结构体名称的复数形式映射生成表名
	db.SingularTable(true)
	//设置表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tb_" + defaultTableName
	}

	return db
}

func InitTable() {
	CreateTalbe(User{})
	CreateTalbe(Interface{})
	CreateTalbe(InterfaceActor{})
	CreateTalbe(Actor{})
}

func CreateTalbe(v interface{}) {
	msdb := GetMysqlDb()
	defer msdb.Close()
	//判断表是否存在，不存在则创建
	if !msdb.HasTable(v) {
		if err := msdb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(v).Error; err != nil {
			panic(err)
		}
	}
}
