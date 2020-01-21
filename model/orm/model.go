package orm

import (
	"github.com/Jarnpher553/micro-core/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

// ModelInt int主键
type ModelInt struct {
	ID          int `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedTime time.Time
	UpdatedTime time.Time
	IsActive    bool `gorm:"not null"`
}

// ModelUUID uuid主键
type ModelUUID struct {
	ID          uuid.GUID `gorm:"primary_key"`
	CreatedTime time.Time `gorm:"not null"`
	UpdatedTime time.Time `gorm:"not null"`
	IsActive    bool      `gorm:"not null"`
}

// BeforeCreate 自动主键值
func (m ModelUUID) BeforeCreate(scope *gorm.Scope) error {
	field, _ := scope.FieldByName("ID")
	if field.IsBlank {
		if err := scope.SetColumn("ID", uuid.New()); err != nil {
			return err
		}
	}

	if err := scope.SetColumn("CreatedTime", time.Now()); err != nil {
		return err
	}
	if err := scope.SetColumn("UpdatedTime", time.Now()); err != nil {
		return err
	}
	if err := scope.SetColumn("IsActive", true); err != nil {
		return err
	}
	return nil
}

func (m ModelUUID) BeforeUpdate(scope *gorm.Scope) (err error) {
	if err := scope.SetColumn("UpdatedTime", time.Now()); err != nil {
		return err
	}
	return nil
}

func (m ModelUUID) BeforeDelete(scope *gorm.Scope) (err error) {
	if err := scope.SetColumn("IsActive", false); err != nil {
		return err
	}
	return nil
}

// BeforeCreate 自动主键值
func (m ModelInt) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("CreatedTime", time.Now()); err != nil {
		return err
	}
	if err := scope.SetColumn("UpdatedTime", time.Now()); err != nil {
		return err
	}
	if err := scope.SetColumn("IsActive", true); err != nil {
		return err
	}
	return nil
}

func (m ModelInt) BeforeUpdate(scope *gorm.Scope) (err error) {
	if err := scope.SetColumn("UpdatedTime", time.Now()); err != nil {
		return err
	}
	return nil
}

func (m ModelInt) BeforeDelete(scope *gorm.Scope) (err error) {
	if err := scope.SetColumn("IsActive", false); err != nil {
		return err
	}
	return nil
}
