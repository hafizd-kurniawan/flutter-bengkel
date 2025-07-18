package services

import (
	"errors"
	"time"

	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/models"
	"flutter-bengkel/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (*models.LoginResponse, error)
	RefreshToken(refreshToken string) (*models.LoginResponse, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	GenerateTokens(user *models.User) (string, string, error)
	ValidateToken(tokenString string) (*models.Claims, error)
}

type authService struct {
	repos *repositories.Repositories
	cfg   *config.Config
}

func NewAuthService(repos *repositories.Repositories, cfg *config.Config) AuthService {
	return &authService{
		repos: repos,
		cfg:   cfg,
	}
}

func (s *authService) Login(username, password string) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.repos.User.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password
	if err := s.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := s.repos.User.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail login
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.ExpireHours * 3600, // Convert hours to seconds
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	// Parse and validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user from database
	user, err := s.repos.User.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.cfg.JWT.ExpireHours * 3600,
	}, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *authService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) GenerateTokens(user *models.User) (string, string, error) {
	// Get user permissions
	permissions, err := s.repos.Role.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return "", "", err
	}

	// Convert permissions to string array
	permissionNames := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionNames[i] = permission.Name
	}

	// Create claims for access token
	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"username":    user.Username,
		"role_id":     user.RoleID,
		"outlet_id":   user.OutletID,
		"permissions": permissionNames,
		"exp":         time.Now().Add(time.Hour * time.Duration(s.cfg.JWT.ExpireHours)).Unix(),
		"iat":         time.Now().Unix(),
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	// Create claims for refresh token (longer expiry, minimal data)
	refreshClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.cfg.JWT.RefreshExpireHours)).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.Claims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Create models.Claims
	userClaims := &models.Claims{
		UserID:   int64(claims["user_id"].(float64)),
		Username: claims["username"].(string),
		RoleID:   int64(claims["role_id"].(float64)),
	}

	// Outlet ID is optional
	if outletID, exists := claims["outlet_id"]; exists && outletID != nil {
		id := int64(outletID.(float64))
		userClaims.OutletID = &id
	}

	// Permissions are optional (only in access tokens)
	if permissions, exists := claims["permissions"]; exists {
		if permissionList, ok := permissions.([]interface{}); ok {
			permissionStrings := make([]string, len(permissionList))
			for i, p := range permissionList {
				permissionStrings[i] = p.(string)
			}
			userClaims.Permissions = permissionStrings
		}
	}

	return userClaims, nil
}

// User Service
type UserService interface {
	Create(req *models.CreateUserRequest) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id int64) error
	List(page, limit int) ([]models.User, *models.PaginationMeta, error)
	ChangePassword(id int64, req *models.ChangePasswordRequest) error
}

type userService struct {
	repos *repositories.Repositories
}

func NewUserService(repos *repositories.Repositories) UserService {
	return &userService{repos: repos}
}

func (s *userService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Check if username already exists
	if _, err := s.repos.User.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.repos.User.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		OutletID:     req.OutletID,
		IsActive:     true,
	}

	if err := s.repos.User.Create(user); err != nil {
		return nil, err
	}

	// Get created user with relations
	return s.repos.User.GetByID(user.ID)
}

func (s *userService) GetByID(id int64) (*models.User, error) {
	return s.repos.User.GetByID(id)
}

func (s *userService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	existingUser, err := s.repos.User.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if username is being changed and already exists
	if req.Username != existingUser.Username {
		if _, err := s.repos.User.GetByUsername(req.Username); err == nil {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is being changed and already exists
	if req.Email != existingUser.Email {
		if _, err := s.repos.User.GetByEmail(req.Email); err == nil {
			return nil, errors.New("email already exists")
		}
	}

	// Update user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		OutletID: req.OutletID,
		IsActive: req.IsActive,
	}

	if err := s.repos.User.Update(id, user); err != nil {
		return nil, err
	}

	// Get updated user with relations
	return s.repos.User.GetByID(id)
}

func (s *userService) Delete(id int64) error {
	return s.repos.User.Delete(id)
}

func (s *userService) List(page, limit int) ([]models.User, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	users, total, err := s.repos.User.List(offset, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return users, meta, nil
}

func (s *userService) ChangePassword(id int64, req *models.ChangePasswordRequest) error {
	// Get user
	user, err := s.repos.User.GetByID(id)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.repos.User.ChangePassword(id, string(hashedPassword))
}