package orm

import (
	"github.com/Jarnpher553/gemini/pkg/uuid"
	"time"
)

// ModelInt int主键
type ModelInt struct {
	ID          int       `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	IsActive    bool      `gorm:"not null"`
}

// ModelUUID uuid主键
type ModelUUID struct {
	ID          uuid.GUID `gorm:"primary_key"`
	CreatedTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	IsActive    bool      `gorm:"not null"`
}
