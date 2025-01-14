package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string    `json:"name" form:"name"`
	Email     string    `json:"email" form:"email"`
	Password  string    `json:"password" form:"password"`
	Role      string    `json:"role" form:"role"`
	Dosen     Dosen     `gorm:"foreignKey:user_id"`
	Mahasiswa Mahasiswa `gorm:"foreignKey:user_id"`
}

type Users []User
