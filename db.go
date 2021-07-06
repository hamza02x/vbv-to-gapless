package main

import (
	"strconv"

	hel "github.com/hamza02x/go-helper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var errDB error

type Timing struct {
	Sura int   `gorm:"column:sura;primaryKey"`
	Ayah int   `gorm:"column:ayah;primaryKey"`
	Time int64 `gorm:"column:time"`
}

type Properties struct {
	Property string `gorm:"column:property;primaryKey"`
	Value    string `gorm:"column:value"`
}

func (Timing) TableName() string     { return "timings" }
func (Properties) TableName() string { return "properties" }

func setDB(path string) {

	hel.FileRemoveIfExists(path)

	db, errDB = gorm.Open(sqlite.Open(path), &gorm.Config{})

	panics("Error opening database!", errDB)

	db.AutoMigrate(&Timing{})
	db.AutoMigrate(&Properties{})

	db.Create(&Properties{Property: "version", Value: strconv.Itoa(DB_VERSION)})

	for sura := 1; sura <= TOTAL_SURA; sura++ {
		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
			dbCreateTiming(sura, aya, 0)
		}
		dbCreateTiming(sura, 999, 0)
	}
}

func dbCreateTiming(sura, aya int, time int64) {
	db.Create(&Timing{Sura: sura, Ayah: aya, Time: time})
}

func dbUpdateTiming(sura, aya int, time int64) {
	db.Save(&Timing{Sura: sura, Ayah: aya, Time: time})
}
func dbVaccum() {
	db.Exec("VACUUM")
}
