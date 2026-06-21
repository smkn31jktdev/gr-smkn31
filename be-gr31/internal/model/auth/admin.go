package auth

import "be-gr31/internal/model/common"

type Admin struct {
	ID         string          `json:"id" bson:"_id"`
	Nama       string          `json:"nama" bson:"nama"`
	Email      string          `json:"email" bson:"email"`
	Role       string          `json:"role" bson:"role"`
	IsWalas    bool            `json:"isWalas" bson:"is_walas"`
	Kelas      string          `json:"kelas" bson:"kelas"`
	FotoProfil string          `json:"fotoProfil,omitempty" bson:"fotoProfil,omitempty"`
	Password   string          `json:"password,omitempty" bson:"password"` // hash bcrypt — dikosongkan sebelum response
	CreatedAt  common.FlexTime `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  common.FlexTime `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login Admin
type AdminLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Admin        Admin  `json:"admin"`
}

// Refresh Token Admin
type AdminRefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Tambah Admin
type AdminCreateRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	IsWalas  bool   `json:"isWalas"`
	Kelas    string `json:"kelas"`
	Password string `json:"password" binding:"required,min=8"`
}

// Update Admin
type AdminUpdateRequest struct {
	Role    *string `json:"role"`
	IsWalas *bool   `json:"isWalas"`
	Kelas   *string `json:"kelas"`
}

// List Admin
type AdminFilter struct {
	Role  string
	Query string
	Page  int
	Limit int
}

// Import Bulk Admin
type AdminBulkItem struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password"`
}

// Claim JWT
type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	NISN  string `json:"nisn,omitempty"`
}
