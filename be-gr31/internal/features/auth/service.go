package auth

import (
	"context"
	"errors"
	"strings"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/util"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("NIS atau password salah")
	ErrNotFound           = errors.New("data tidak ditemukan")
	ErrUnauthorized       = errors.New("tidak diizinkan")
)

type Service struct {
	repo      *Repo
	jwtConfig util.JWTConfig
	superAdminEmails []string
}

func NewService(repo *Repo, jwtConfig util.JWTConfig, superAdminEmails []string) *Service {
	return &Service{
		repo:             repo,
		jwtConfig:        jwtConfig,
		superAdminEmails: superAdminEmails,
	}
}

func (s *Service) LoginSiswa(ctx context.Context, req authmodel.SiswaLoginRequest) (*authmodel.SiswaLoginResponse, error) {
	siswa, err := s.repo.FindSiswaByNIS(ctx, req.NIS)
	if err != nil {
		return nil, err
	}
	if siswa == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(siswa.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := util.JWTClaims{
		ID:   siswa.ID,
		Role: "student",
		NIS:  siswa.NIS,
	}

	accessToken, err := util.IssueAccessToken(s.jwtConfig, claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := util.IssueRefreshToken(s.jwtConfig, claims)
	if err != nil {
		return nil, err
	}

	siswa.Password = ""
	return &authmodel.SiswaLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Siswa:        *siswa,
	}, nil
}

func (s *Service) LoginAdmin(ctx context.Context, req authmodel.AdminLoginRequest) (*authmodel.AdminLoginResponse, error) {
	admin, err := s.repo.FindAdminByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	role := admin.Role
	if strings.ToLower(strings.TrimSpace(admin.Email)) == "smkn31jktdev@gmail.com" {
		role = "super_admin"
		if admin.Role != "super_admin" {
			_ = s.repo.UpdateAdminRole(ctx, admin.ID, "super_admin")
			admin.Role = "super_admin"
		}
	}
	if role == "" {
		role = "admin"
	}

	claims := util.JWTClaims{
		ID:    admin.ID,
		Email: admin.Email,
		Role:  role,
	}

	accessToken, err := util.IssueAccessToken(s.jwtConfig, claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := util.IssueRefreshToken(s.jwtConfig, claims)
	if err != nil {
		return nil, err
	}

	admin.Password = ""
	admin.Role = role
	return &authmodel.AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        *admin,
	}, nil
}

func (s *Service) RefreshSiswaToken(ctx context.Context, req authmodel.SiswaRefreshRequest) (*authmodel.SiswaLoginResponse, error) {
	claims, err := util.ParseJWT(req.RefreshToken, s.jwtConfig.Secret)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if claims.Role != "student" {
		return nil, ErrUnauthorized
	}

	siswa, err := s.repo.FindSiswaByID(ctx, claims.ID)
	if err != nil || siswa == nil {
		return nil, ErrNotFound
	}

	newClaims := util.JWTClaims{
		ID:   siswa.ID,
		Role: "student",
		NIS:  siswa.NIS,
	}
	accessToken, err := util.IssueAccessToken(s.jwtConfig, newClaims)
	if err != nil {
		return nil, err
	}
	newRefresh, err := util.IssueRefreshToken(s.jwtConfig, newClaims)
	if err != nil {
		return nil, err
	}

	siswa.Password = ""
	return &authmodel.SiswaLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		Siswa:        *siswa,
	}, nil
}

func (s *Service) RefreshAdminToken(ctx context.Context, req authmodel.AdminRefreshRequest) (*authmodel.AdminLoginResponse, error) {
	claims, err := util.ParseJWT(req.RefreshToken, s.jwtConfig.Secret)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if !isAdminRoleStr(claims.Role) {
		return nil, ErrUnauthorized
	}

	admin, err := s.repo.FindAdminByID(ctx, claims.ID)
	if err != nil || admin == nil {
		return nil, ErrNotFound
	}

	role := admin.Role
	if strings.ToLower(strings.TrimSpace(admin.Email)) == "smkn31jktdev@gmail.com" {
		role = "super_admin"
	}
	if role == "" {
		role = "admin"
	}

	newClaims := util.JWTClaims{
		ID:    admin.ID,
		Email: admin.Email,
		Role:  role,
	}
	accessToken, err := util.IssueAccessToken(s.jwtConfig, newClaims)
	if err != nil {
		return nil, err
	}
	newRefresh, err := util.IssueRefreshToken(s.jwtConfig, newClaims)
	if err != nil {
		return nil, err
	}

	admin.Password = ""
	admin.Role = role
	return &authmodel.AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		Admin:        *admin,
	}, nil
}

func (s *Service) MeSiswa(ctx context.Context, id string) (*authmodel.Siswa, error) {
	siswa, err := s.repo.FindSiswaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if siswa == nil {
		return nil, ErrNotFound
	}
	siswa.Password = ""
	return siswa, nil
}

func (s *Service) MeAdmin(ctx context.Context, id string) (*authmodel.Admin, error) {
	admin, err := s.repo.FindAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrNotFound
	}
	admin.Password = ""
	if strings.ToLower(strings.TrimSpace(admin.Email)) == "smkn31jktdev@gmail.com" {
		admin.Role = "super_admin"
	}
	if admin.Role == "" {
		admin.Role = "admin"
	}
	return admin, nil
}

func (s *Service) isSuperAdmin(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, e := range s.superAdminEmails {
		if strings.ToLower(strings.TrimSpace(e)) == email {
			return true
		}
	}
	return false
}

func isAdminRoleStr(role string) bool {
	switch role {
	case "admin", "super_admin", "guru_bk", "guru_wali", "piket", "bk", "admin_bk", "walas", "admin_piket":
		return true
	}
	return false
}
