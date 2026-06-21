package auth

import (
	"errors"
	"log"
	"net/http"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Handler autentikasi
type Handler struct {
	service *Service
}

// Membuat instance handler baru
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// LoginSiswa
func (h *Handler) LoginSiswa(c *gin.Context) {
	var req authmodel.SiswaLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.LoginSiswa(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, common.Fail(err.Error()))
			return
		}
		log.Printf("ERROR LoginSiswa: %v", err)
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "login berhasil"))
}

// MeSiswa
func (h *Handler) MeSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)
	siswa, err := h.service.MeSiswa(c.Request.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("siswa tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}
	c.JSON(http.StatusOK, common.OK(siswa, "ok"))
}

// RefreshTokenSiswa
func (h *Handler) RefreshTokenSiswa(c *gin.Context) {
	var req authmodel.SiswaRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.RefreshSiswaToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, common.Fail("refresh token tidak valid"))
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("siswa tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "token diperbarui"))
}

// LoginAdmin
func (h *Handler) LoginAdmin(c *gin.Context) {
	var req authmodel.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.LoginAdmin(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "login berhasil"))
}

// MeAdmin
func (h *Handler) MeAdmin(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)
	admin, err := h.service.MeAdmin(c.Request.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("admin tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}
	c.JSON(http.StatusOK, common.OK(admin, "ok"))
}

// RefreshTokenAdmin
func (h *Handler) RefreshTokenAdmin(c *gin.Context) {
	var req authmodel.AdminRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.RefreshAdminToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, common.Fail("refresh token tidak valid"))
			return
		}
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("admin tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "token diperbarui"))
}
