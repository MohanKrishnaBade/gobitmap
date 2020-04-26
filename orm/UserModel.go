package orm

import "github.com/jinzhu/gorm"

// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` into model `User`
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
	Password string `json:"password"`
}

func (u *User) CreateTable(db *gorm.DB) {
	if !db.HasTable(u.tableName()) {
		tx := db.Begin()
		err := tx.Table(u.tableName()).CreateTable(&u).Error
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}
}
func (u *User) Create(db *gorm.DB) {

	tx := db.Begin()
	err := tx.Table(u.tableName()).Create(&u).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func (u *User) Update(db *gorm.DB) {
	db.Table(u.tableName()).Omit("password").Save(&u)
}

func (u *User) tableName() string {
	return "user"
}

func (u *User) Delete(db *gorm.DB) {
	tx := db.Begin()
	err := tx.Table(u.tableName()).Delete(&u).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()

}
func (u *User) SearchById(db *gorm.DB) {
	tx := db.Begin()
	err := tx.Table(u.tableName()).Where(map[string]interface{}{"id": u.ID}).First(&u).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func (u *User) SearchByEmail(db *gorm.DB) {
	tx := db.Begin()
	err := tx.Table(u.tableName()).Where(map[string]interface{}{"email": u.Email}).First(&u).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}
