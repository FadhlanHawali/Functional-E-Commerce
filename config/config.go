package config

import (
	"github.com/jinzhu/gorm"
)

func DBInit() *gorm.DB {
	//TODO CONFIG DB
	db, err := gorm.Open("mysql", "admin:OGTIRKWUIGESPYPB@tcp(sl-aus-syd-1-portal.5.dblayer.com:20533)/piruma")
	//db, err := gorm.Open("postgres", "host=localhost port=4000 user=mydatabase dbname=piruma password=mydatabase sslmode=disable")
	if err != nil {
		panic("failed to connect to database" + err.Error())
	}

	//TODO AUTO MIGRATE SCHEMA
	//db.AutoMigrate(model.Mahasiswa{},model.History{},model.Bank{})
	return db
}
