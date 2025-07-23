// migrations/migrations.go
package migrations

import (
	"your_project/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM AutoMigrate for all models
func AutoMigrate(db *gorm.DB) error {
	// Add all your models here for auto-migration
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	// Add other models here as needed:
	// err = db.AutoMigrate(&model.AnotherModel{})
	// if err != nil {
	//	return err
	// }

	return nil
}

// You can define more complex migrations here if needed,
// for example, using raw SQL or GORM's migration features
// func CustomMigration1(db *gorm.DB) error {
//	// Perform custom migration steps
//	return nil
// }
