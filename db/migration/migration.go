package migration

import (
	"Test_REST/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "19122024",
			Migrate: func(db *gorm.DB) error {
				var w models.Wallet
				return db.Migrator().CreateTable(&w)
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("wallet")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Migration success")
}
