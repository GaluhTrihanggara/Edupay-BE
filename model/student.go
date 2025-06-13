package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Student adalah model untuk data siswa
type Student struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk siswa
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`                       // Nama depan siswa
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`                        // Nama belakang siswa
	Class     string         `gorm:"type:varchar(50)" json:"class"`                             // Kelas siswa
	UserID    uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`                            // Foreign key untuk User (orang tua)
	User      User           `gorm:"foreignKey:UserID" json:"user"`                             // Relasi ke model User
	CreatedAt time.Time      `json:"created_at"`                                                // Timestamp pembuatan
	UpdatedAt time.Time      `json:"updated_at"`                                                // Timestamp pembaruan
	DeletedAt gorm.DeletedAt `json:"deleted_at"`                                                // Timestamp penghapusan

}
