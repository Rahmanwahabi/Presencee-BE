package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string    `json:"name" form:"name"`
	Email     string    `json:"email" form:"email"`
	Password  string    `json:"password" form:"password"`
	Token     string    `gorm:"-"`
	Role      string    `json:"role" form:"role"`
	Image     string    `json:"image" form:"image"`
	Dosen     Dosen     `gorm:"foreignKey:user_id"`
	Mahasiswa Mahasiswa `gorm:"foreignKey:user_id"`
}

type Mahasiswa struct {
	gorm.Model
	Mahasiswa  string `json:"mahasiswa" form:"mahasiswa"`
	Email      string `json:"email" form:"email"`
	Name       string `json:"name" form:"name"`
	NIM        string `json:"nim" form:"nim"`
	Image      string `json:"image" form:"image"`
	Phone      string `json:"phone" form:"phone"`
	Jurusan    string `json:"jurusan" form:"jurusan"`
	TahunMasuk string `json:"tahun_masuk" form:"tahun_masuk"`
	IPK        string `json:"ipk" form:"ipk"`
	UserID     uint   `json:"user_id" form:"user_id"`
}

type Dosen struct {
	gorm.Model
	Nama   string `json:"nama" form:"nama"`
	Email  string `json:"email" form:"email"`
	NIP    string `json:"nip" form:"nip"`
	Image  string `json:"image" form:"image"`
	UserID uint   `json:"user_id" form:"user_id"`
}
