package main

import (
	"strconv"

	hel "github.com/hamza02x/go-helper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var errDB error

type TimingUnordered Timing

type Timing struct {
	Sura int   `gorm:"column:sura;primaryKey"`
	Ayah int   `gorm:"column:ayah;primaryKey"`
	Time int64 `gorm:"column:time"`
}

type Properties struct {
	Property string `gorm:"column:property;primaryKey"`
	Value    string `gorm:"column:value"`
}

func (Timing) TableName() string          { return "timings" }
func (TimingUnordered) TableName() string { return "timings_unordered" }
func (Properties) TableName() string      { return "properties" }

func setDB(path string) {

	hel.FileRemoveIfExists(path)

	db, errDB = gorm.Open(sqlite.Open(path), &gorm.Config{})

	panics("Error opening database!", errDB)

	db.AutoMigrate(&TimingUnordered{})
	db.AutoMigrate(&Timing{})
	db.AutoMigrate(&Properties{})

	db.Create(&Properties{Property: "version", Value: strconv.Itoa(DB_VERSION)})
}
