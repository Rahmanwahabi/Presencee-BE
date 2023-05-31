package payload

type (
	CreateUserRequest struct {
		Name     string `json:"name" form:"name" validate:"required,max=20"`
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required,min=5"`
	}

	CreateUserResponse struct {
		UserID uint   `json:"user_id"`
		Token  string `json:"token"`
	}

	LoginUserRequest struct {
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required,min=5"`
	}

	LoginUserResponse struct {
		UserID uint   `json:"user_id"`
		Token  string `json:"token"`
	}
)

type CreateDosenRequest struct {
	Nama   string `json:"nama" form:"nama" validate:"required"`
	Email  string `json:"email" form:"email" validate:"required,email"`
	NIP    string `json:"nip" form:"nip" validate:"required"`
	Image  string `json:"image" form:"image"`
	UserID uint   `json:"user_id" form:"user_id" validate:"required"`
}

type CreateDosenResponse struct {
	DosenID uint `json:"dosen_id"`
}

type UpdateDosenRequest struct {
	Nama   string `json:"nama" form:"nama" validate:"required"`
	Email  string `json:"email" form:"email" validate:"required,email"`
	NIP    string `json:"nip" form:"nip" validate:"required"`
	Image  string `json:"image" form:"image"`
	UserID uint   `json:"user_id" form:"user_id" validate:"required"`
}

type UpdateDosenResponse struct {
	DosenID uint `json:"dosen_id"`
}

type GetDosenResponse struct {
	DosenID uint   `json:"dosen_id"`
	Nama    string `json:"nama"`
	Email   string `json:"email"`
	NIP     string `json:"nip"`
	Image   string `json:"image"`
	UserID  uint   `json:"user_id"`
}

type GetDosensResponse struct {
	Dosens []GetDosenResponse `json:"dosens"`
}
