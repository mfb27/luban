package admin

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mfb27/luban/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrAdminNotFound    = errors.New("admin not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrAdminDisabled    = errors.New("admin account is disabled")
)

// AdminClaims JWT claims for admin
type AdminClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AdminAuthService 管理员认证服务
type AdminAuthService struct {
	db     *gorm.DB
	secret string
}

// NewAdminAuthService 创建管理员认证服务
func NewAdminAuthService(db *gorm.DB, secret string) *AdminAuthService {
	return &AdminAuthService{
		db:     db,
		secret: secret,
	}
}

// Login 管理员登录
func (s *AdminAuthService) Login(email, password string) (string, error) {
	var admin model.Admin
	if err := s.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ErrAdminNotFound
		}
		return "", err
	}

	// 检查管理员状态
	if admin.Status != "active" {
		return "", ErrAdminDisabled
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(&admin).Update("updated_at", now)

	// 生成JWT token
	claims := AdminClaims{
		UserID: admin.ID,
		Email:  admin.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "luban",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ValidateToken 验证JWT token
func (s *AdminAuthService) ValidateToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// CreateAdmin 创建管理员（初始管理员）
func (s *AdminAuthService) CreateAdmin(name, email, password string) error {
	// 检查是否已存在管理员
	var count int64
	s.db.Model(&model.Admin{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New("admin already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.Admin{
		ID:       generateID(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Status:   "active",
	}

	return s.db.Create(&admin).Error
}

// GetAdminByID 根据ID获取管理员
func (s *AdminAuthService) GetAdminByID(id string) (*model.Admin, error) {
	var admin model.Admin
	err := s.db.First(&admin, "id = ?", id).Error
	return &admin, err
}

// GetAdminByEmail 根据邮箱获取管理员
func (s *AdminAuthService) GetAdminByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	err := s.db.First(&admin, "email = ?", email).Error
	return &admin, err
}

// generateID 生成ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "admin_" + hex.EncodeToString(bytes)
}