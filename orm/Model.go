package orm

import "github.com/jinzhu/gorm"

type Model interface {
	CreateTable(db *gorm.DB)
	Create(db *gorm.DB)
	Delete(db *gorm.DB)
	tableName() string
	SearchByEmail(db *gorm.DB)
	SearchById(db *gorm.DB)
	Update(db *gorm.DB)
}
