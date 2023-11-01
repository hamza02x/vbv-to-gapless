package main

import (
	"strconv"

	hel "github.com/hamza72x/go-helper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

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

func setDB(path string) error {

	if err := hel.FileRemoveIfExists(path); err != nil {
		return err
	}

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&Timing{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Properties{}); err != nil {
		return err
	}

	if err := db.Create(&Properties{Property: "version", Value: strconv.Itoa(DB_VERSION)}).Error; err != nil {
		return err
	}

	for sura := 1; sura <= TOTAL_SURA; sura++ {
		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
			if err := dbCreateTiming(sura, aya, 0); err != nil {
				return err
			}
		}
		if err := dbCreateTiming(sura, 999, 0); err != nil {
			return err
		}
	}

	return nil
}

func dbCreateTiming(sura, aya int, time int64) error {
	return db.Create(&Timing{Sura: sura, Ayah: aya, Time: time}).Error
}

func dbUpdateTiming(sura, aya int, time int64) error {
	return db.Save(&Timing{Sura: sura, Ayah: aya, Time: time}).Error
}

func dbVaccum() error {
	return db.Exec("VACUUM").Error
}
