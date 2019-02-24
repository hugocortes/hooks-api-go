package migrations

import (
	"github.com/hugocortes/hooks-api/bins/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	gormigrate "gopkg.in/gormigrate.v1"
)

var m *gormigrate.Gormigrate

// Run initializes the schema and runs any additional migrations
func Run(db *gorm.DB) {
	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// migrations
	})

	initSchema()
	migrate()
}

func initSchema() {
	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&models.Bin{},
		).Error
		if err != nil {
			return err
		}

		logrus.Debug("Initialize schema √")
		return nil
	})
}

func migrate() {
	if err := m.Migrate(); err != nil {
		logrus.Fatal("Could not migrate: ", err)
	}
	logrus.Debug("Migration √")
}
