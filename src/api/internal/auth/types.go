package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secret123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Role represents a user role
type Role struct {
	ID           string    `json:"id" example:"1"`
	RoleID       int       `json:"role_id"`
	RoleName     string    `json:"role_name"`
	RoleDesc     string    `json:"role_desc,omitempty"`
	Name         string    `json:"name" example:"admin"`
	Permissions  []string  `json:"permissions" example:"read,write"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	StatusID     int       `json:"status_id"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" example:"1"`
	UserID    int       `json:"user_id"` // Keep for backward compatibility
	Username  string    `json:"username" example:"johndoe"`
	Password  string    `json:"password,omitempty"`
	Salt      string    `json:"salt,omitempty"`
	Email     string    `json:"email" example:"john@example.com"`
	RoleID    string    `json:"role_id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	StatusID  int       `json:"status_id"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secret123"`
	Email    string `json:"email" example:"john@example.com"`
}

type UserRoles struct {
	UserRoleID int       `json:"user_role_id"`
	RoleID     int       `json:"role_id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	StatusID   int       `json:"status_id"`
}

type RolePermissions struct {
	RolePermissionID   int       `json:"role_permission_id"`
	RolePermissionDesc string    `json:"role_permission_desc"`
	ResourceTypeID     int       `json:"resource_type_id"`
	ResourceName       string    `json:"resource_name"`
	CanExecute         bool      `json:"can_execute"`
	CanRead            bool      `json:"can_read"`
	CanWrite           bool      `json:"can_write"`
	CanDelete          bool      `json:"can_delete"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	StatusID           int       `json:"status_id"`
}

type PermissionView struct {
	UserID       string
	RoleID       int
	RoleName     string
	ResourceName string
	CanExecute   bool
	CanRead      bool
	CanWrite     bool
	CanDelete    bool
}

type JwtClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}

type JwtToken struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}

type UserPermissionView struct {
	UserID           string
	RoleID           int
	RoleName         string
	IsSuperAdmin     bool
	RolePermissionID int
	ResourceTypeID   int
	ResourceName     string
	CanExecute       bool
	CanRead          bool
	CanWrite         bool
	CanDelete        bool
}

type IsSuperAdmin struct {
	UserID       int
	IsSuperAdmin int
}
