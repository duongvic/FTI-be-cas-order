package migrations

import (
	"casorder/db/migrations/models"
	"casorder/db/migrations/versions"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func Upgrate(db *gorm.DB) error {
	return versions.Update(db)
}

func Seed(db *gorm.DB) {
	fmt.Println("\n========================= Seeding Database =========================")

	fmt.Println("Seeding Regions.....")
	if err := db.First(&models.Region{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		for _, r := range RegionData {
			db.Model(&r).Create(&r)
		}
	}

	fmt.Println("Seeding Units.....")
	if err := db.First(&models.Unit{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		for _, u := range UnitData {
			db.Model(&u).Create(&u)
		}
	}

	fmt.Println("Seeding Products.....")
	if err := db.First(&models.Product{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		for _, p := range ProductData {
			db.Model(&p).Create(&p)
		}
	}

	fmt.Println("Seeding Products.....")
	if err := db.First(&models.Package{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		for _, pkg := range PackageData {
			db.Model(&pkg).Create(&pkg)
		}
	}

	fmt.Println("Seeding PackageProducts.....")
	if err := db.First(&models.PackageProduct{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		for _, pp := range PackageProductData {
			db.Model(&pp).Create(&pp)
		}
	}

	fmt.Println("========================= Seeding Complete! =========================")
}