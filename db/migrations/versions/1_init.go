package versions

import (
	"gorm.io/gorm"

	"casorder/db/migrations/models"
)

func Update(db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		new(models.Order),
		new(models.Product),
		new(models.Unit),
		new(models.PackageProduct),
		new(models.OrderProduct),
		new(models.OrderDtl),
		new(models.Package),
		new(models.Region),
		new(models.Task),
		new(models.Transaction),
		new(models.User),
	)
}
